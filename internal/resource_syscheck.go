package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSyscheck manages Wazuh file integrity monitoring (FIM/Syscheck)
// for a specific agent via:
// - PUT    /syscheck            (run scan with agents_list)
// - GET    /syscheck/{agent_id} (get findings summary)
// - DELETE /syscheck/{agent_id} (clear results, only for older agents)
func resourceSyscheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSyscheckCreate,
		ReadContext:   resourceSyscheckRead,
		// No UpdateContext (all fields are ForceNew or Computed)
		DeleteContext: resourceSyscheckDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSyscheckImport,
		},

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh agent ID (e.g. '001') for which syscheck (FIM) is managed.",
			},

			// --- Scan result (PUT /syscheck) ---
			"scan_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh when the syscheck scan is started.",
			},
			"scan_total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents for which the syscheck scan was restarted.",
			},
			"scan_total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents where starting the syscheck scan failed.",
			},

			// --- Results summary (GET /syscheck/{agent_id}) ---
			"results_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned when fetching FIM findings for the agent.",
			},
			"results_total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of FIM findings (items) returned for this agent.",
			},
			"results_total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of failed items when retrieving FIM findings.",
			},

			"last_scan_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the last syscheck scan was triggered via Terraform.",
			},
		},
	}
}

// Create: run FIM scan for the agent via PUT /syscheck?agents_list=agent_id
func resourceSyscheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Get("agent_id").(string)

	// Build URL: /syscheck?agents_list=001
	urlStr := fmt.Sprintf("%s/syscheck", client.Endpoint)
	query := url.Values{}
	query.Set("agents_list", agentID)
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	// Body can be an empty JSON object
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to start syscheck scan for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			TotalAffected int `json:"total_affected_items"`
			TotalFailed   int `json:"total_failed_items"`
		} `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}

	_ = json.Unmarshal(respBody, &result)

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error while starting syscheck scan: %s", string(respBody))
	}

	_ = d.Set("scan_message", result.Message)
	_ = d.Set("scan_total_affected", result.Data.TotalAffected)
	_ = d.Set("scan_total_failed", result.Data.TotalFailed)
	_ = d.Set("last_scan_timestamp", time.Now().UTC().Format(time.RFC3339))

	// agent_id is the resource ID
	d.SetId(agentID)

	return diags
}

// Read: get FIM findings summary via GET /syscheck/{agent_id}
func resourceSyscheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Id()
	urlStr := fmt.Sprintf("%s/syscheck/%s", client.Endpoint, agentID)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		// No syscheck data for this agent
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read syscheck results for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			TotalAffected int `json:"total_affected_items"`
			TotalFailed   int `json:"total_failed_items"`
		} `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse syscheck results for agent '%s': %v", agentID, err)
	}

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error when getting syscheck results: %s", string(body))
	}

	_ = d.Set("results_message", result.Message)
	_ = d.Set("results_total_affected", result.Data.TotalAffected)
	_ = d.Set("results_total_failed", result.Data.TotalFailed)
	_ = d.Set("agent_id", agentID)

	return diags
}

// Delete: clear FIM scan results for the agent via DELETE /syscheck/{agent_id}
// (only applies for agents < 3.12.0; for newer agents, API may be a no-op but should still return success)
func resourceSyscheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Id()
	urlStr := fmt.Sprintf("%s/syscheck/%s", client.Endpoint, agentID)

	req, err := http.NewRequest("DELETE", urlStr, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Treat 404 as "already cleared"
	if resp.StatusCode != http.StatusNotFound && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return diag.Errorf("failed to clear syscheck database for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import: attach existing agent_id as syscheck-managed resource (no API call)
func resourceSyscheckImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	agentID := d.Id()
	if err := d.Set("agent_id", agentID); err != nil {
		return nil, fmt.Errorf("failed to set agent_id during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
