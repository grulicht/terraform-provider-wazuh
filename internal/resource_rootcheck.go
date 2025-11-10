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

// resourceRootcheck manages Wazuh rootcheck for a specific agent via:
// - PUT /rootcheck            (run scan with agents_list)
// - GET /rootcheck/{agent_id} (get results)
// - DELETE /rootcheck/{agent_id} (clear results)
func resourceRootcheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRootcheckCreate,
		ReadContext:   resourceRootcheckRead,
		// UpdateContext: resourceRootcheckUpdate,
		DeleteContext: resourceRootcheckDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRootcheckImport,
		},

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh agent ID (e.g. '001') for which rootcheck is managed.",
			},

			// --- Scan result (PUT /rootcheck) ---
			"scan_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh when starting rootcheck scan.",
			},
			"scan_total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents for which rootcheck scan was restarted (should include this agent).",
			},
			"scan_total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents where starting rootcheck scan failed.",
			},

			// --- Results summary (GET /rootcheck/{agent_id}) ---
			"results_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned when fetching rootcheck results.",
			},
			"results_total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of rootcheck items returned for this agent.",
			},
			"results_total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of failed items when retrieving rootcheck results.",
			},

			"last_scan_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the last rootcheck scan was triggered via Terraform.",
			},
		},
	}
}

// Create: trigger rootcheck scan for the agent via PUT /rootcheck?agents_list=agent_id
func resourceRootcheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Get("agent_id").(string)

	// Build URL: /rootcheck?agents_list=001
	urlStr := fmt.Sprintf("%s/rootcheck", client.Endpoint)
	query := url.Values{}
	query.Set("agents_list", agentID)
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	// Body can be empty JSON object
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
		return diag.Errorf("failed to start rootcheck scan for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(respBody))
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
		return diag.Errorf("Wazuh API returned error while starting rootcheck scan: %s", string(respBody))
	}

	_ = d.Set("scan_message", result.Message)
	_ = d.Set("scan_total_affected", result.Data.TotalAffected)
	_ = d.Set("scan_total_failed", result.Data.TotalFailed)
	_ = d.Set("last_scan_timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use agent_id as resource ID
	d.SetId(agentID)

	return diags
}

// Read: get rootcheck results summary via GET /rootcheck/{agent_id}
func resourceRootcheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Id()
	urlStr := fmt.Sprintf("%s/rootcheck/%s", client.Endpoint, agentID)

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
		// agent rootcheck data no longer exists
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read rootcheck results for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(body))
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
		return diag.Errorf("failed to parse rootcheck results for agent '%s': %v", agentID, err)
	}

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error when getting rootcheck results: %s", string(body))
	}

	_ = d.Set("results_message", result.Message)
	_ = d.Set("results_total_affected", result.Data.TotalAffected)
	_ = d.Set("results_total_failed", result.Data.TotalFailed)

	// keep agent_id in state
	_ = d.Set("agent_id", agentID)

	return diags
}

func resourceRootcheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRootcheckCreate(ctx, d, meta)
}

// Delete: clear agent's rootcheck db via DELETE /rootcheck/{agent_id}
func resourceRootcheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := d.Id()
	urlStr := fmt.Sprintf("%s/rootcheck/%s", client.Endpoint, agentID)

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

	if resp.StatusCode != http.StatusNotFound && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return diag.Errorf("failed to clear rootcheck database for agent '%s': status %d, body: %s", agentID, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

func resourceRootcheckImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	agentID := d.Id()
	if err := d.Set("agent_id", agentID); err != nil {
		return nil, fmt.Errorf("failed to set agent_id during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
