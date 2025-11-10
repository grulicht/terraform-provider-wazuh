package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAgent manages Wazuh agents via:
//   - POST /agents/insert   (Create)
//   - GET  /agents          (Read)
//   - DELETE /agents        (Delete)
func resourceAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentCreate,
		ReadContext:   resourceAgentRead,
		DeleteContext: resourceAgentDelete,

		// Allow importing existing agents by ID:
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgentImport,
		},

		Schema: map[string]*schema.Schema{
			// ---- Inputs / Identity ----

			// Wazuh agent name (required by API)
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Agent name in Wazuh.",
			},

			// Agent ID (optional on create, computed from API if not provided)
			"agent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Wazuh agent ID. If not set, the manager may assign one and it will be populated from the API response.",
			},

			// Agent IP used when registering / communicating with the manager
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IP or IP/NET or ANY. If omitted, Wazuh will try to detect it automatically.",
			},

			// Shared key; optional. If not provided, you typically register agents via other mechanisms.
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Shared key used for communication with the manager. If omitted, Wazuh may generate / manage it separately.",
			},

			// Whether to purge the agent from keystore on destroy
			"purge_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, permanently delete the agent from the key store on destroy (DELETE /agents with purge=true).",
			},

			// ---- Computed fields from GET /agents ----

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent status (e.g. active, pending, never_connected, disconnected).",
			},
			"register_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP used when registering the agent.",
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Manager hostname where the agent is connected.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent version.",
			},
			"node_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster node name the agent is connected to (if any).",
			},
			// Force insertion behaviour (optional)
			"force_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, enable force insertion. Allows replacing existing agents matching name/ID/IP according to additional force conditions.",
			},
			"force_disconnected_time_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Whether to enforce the disconnected_time condition when using force.",
			},
			"force_disconnected_time_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1h",
				ForceNew:    true,
				Description: "Time the agent must have been disconnected to allow forced insertion (e.g. \"30m\", \"2h\", \"7d\").",
			},
			"force_after_registration_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1h",
				ForceNew:    true,
				Description: "Time the agent must have been registered to allow forced insertion (e.g. \"1h\", \"2h\", \"7d\").",
			},
		},
	}
}

// Create: POST /agents/insert
func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	name := strings.TrimSpace(d.Get("name").(string))
	if name == "" {
		return diag.Errorf("name must not be empty")
	}

	agentID := strings.TrimSpace(d.Get("agent_id").(string))
	ip := strings.TrimSpace(d.Get("ip").(string))
	key := strings.TrimSpace(d.Get("key").(string))

	payload := make(map[string]interface{})
	payload["name"] = name
	if agentID != "" {
		payload["id"] = agentID
	}
	if ip != "" {
		payload["ip"] = ip
	}
	if key != "" {
		payload["key"] = key
	}

	// --- Force options ---
	forceEnabled := d.Get("force_enabled").(bool)
	if forceEnabled {
		force := map[string]interface{}{
			"enabled": true,
			"disconnected_time": map[string]interface{}{
				"enabled": d.Get("force_disconnected_time_enabled").(bool),
				"value":   d.Get("force_disconnected_time_value").(string),
			},
			"after_registration_time": d.Get("force_after_registration_time").(string),
		}
		payload["force"] = force
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/agents/insert", client.Endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to create agent '%s': status %d, body: %s", name, resp.StatusCode, string(respBody))
	}

	// Parse response: { "data": { "id": "010", "key": "..." }, "error": 0 }
	var result struct {
		Data struct {
			ID  string `json:"id"`
			Key string `json:"key"`
		} `json:"data"`
		Error int `json:"error"`
	}

	_ = json.Unmarshal(respBody, &result)

	finalID := agentID
	if finalID == "" {
		finalID = result.Data.ID
	}
	if finalID == "" {
		return diag.Errorf("agent created but no ID returned by Wazuh")
	}

	d.SetId(finalID)
	_ = d.Set("agent_id", finalID)

	if result.Data.Key != "" {
		_ = d.Set("key", result.Data.Key)
	}

	// Read back details to populate computed attributes
	return resourceAgentRead(ctx, d, meta)
}

// Read: GET /agents?agents_list=<id>
func resourceAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/agents", client.Endpoint)

	q := url.Values{}
	q.Set("agents_list", id)
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
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

	// 404 or empty list => agent no longer exists
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read agent '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				IP         string `json:"ip"`
				RegisterIP string `json:"registerIP"`
				Status     string `json:"status"`
				Manager    string `json:"manager"`
				Version    string `json:"version"`
				NodeName   string `json:"node_name"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse agent read response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		// Agent not found
		d.SetId("")
		return diags
	}

	item := result.Data.AffectedItems[0]

	_ = d.Set("agent_id", item.ID)
	_ = d.Set("name", item.Name)
	_ = d.Set("ip", item.IP)
	_ = d.Set("register_ip", item.RegisterIP)
	_ = d.Set("status", item.Status)
	_ = d.Set("manager", item.Manager)
	_ = d.Set("version", item.Version)
	_ = d.Set("node_name", item.NodeName)

	return diags
}

// Delete: DELETE /agents?agents_list=<id>&status=all&older_than=0s&purge=<purge>
func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/agents", client.Endpoint)

	purge := d.Get("purge_on_destroy").(bool)

	q := url.Values{}
	q.Set("agents_list", id)
	q.Set("status", "all")    // allow deletion regardless of status
	q.Set("older_than", "0s") // consider all agents
	if purge {
		q.Set("purge", "true")
	} else {
		q.Set("purge", "false")
	}

	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", urlStr, nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to delete agent '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	d.SetId("")
	return diags
}

// Import: terraform import wazuh_agent.example 001
func resourceAgentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// ID passed on import is the Wazuh agent ID
	_ = d.Set("agent_id", d.Id())
	return []*schema.ResourceData{d}, nil
}
