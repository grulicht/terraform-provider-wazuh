package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAgentNodeRestart models an action that restarts all agents
// belonging to a specific cluster node via PUT /agents/node/{node_id}/restart.
func resourceAgentNodeRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentNodeRestartCreate,
		ReadContext:   resourceAgentNodeRestartRead,
		// No Update – one-shot action
		DeleteContext: resourceAgentNodeRestartDelete,

		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster node name whose agents should be restarted.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after sending the restart command.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents for which the restart command was processed.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents where the restart command failed.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw error code returned by Wazuh API (0 = success).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the restart request was sent.",
			},
		},
	}
}

// Create: send restart command via PUT /agents/node/{node_id}/restart
func resourceAgentNodeRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	nodeID := d.Get("node_id").(string)

	urlStr := fmt.Sprintf("%s/agents/node/%s/restart", client.Endpoint, nodeID)

	// Minimal JSON body for compatibility
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to restart agents on node '%s': status %d, body: %s", nodeID, resp.StatusCode, string(respBody))
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

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("error_code", result.Error)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Unique ID = node + timestamp
	id := fmt.Sprintf("%s-%s", nodeID, time.Now().UTC().Format("20060102T150405Z"))
	d.SetId(id)

	return diags
}

// Read: no-op, we don't re-query Wazuh for an action resource
func resourceAgentNodeRestartRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: remove from state only
func resourceAgentNodeRestartDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
