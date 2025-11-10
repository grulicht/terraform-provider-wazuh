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

// resourceAgentRestartGroup models an action that restarts all agents
// belonging to a specific group via PUT /agents/group/{group_id}/restart.
func resourceAgentRestartGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentRestartGroupCreate,
		ReadContext:   resourceAgentRestartGroupRead,
		// One-shot action – no Update
		DeleteContext: resourceAgentRestartGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh group ID (group name) whose agents will be restarted.",
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
				Description: "Raw error code returned by Wazuh API (0 = success, >0 indicates partial/failed states).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the restart request was sent.",
			},
		},
	}
}

func resourceAgentRestartGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Get("group_id").(string)
	if groupID == "" {
		return diag.Errorf("group_id must not be empty")
	}

	urlStr := fmt.Sprintf("%s/agents/group/%s/restart", client.Endpoint, groupID)
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
		return diag.Errorf("failed to restart agents in group '%s': status %d, body: %s", groupID, resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			TotalAffected int    `json:"total_affected_items"`
			TotalFailed   int    `json:"total_failed_items"`
			Message       string `json:"message"`
			Error         int    `json:"error"`
		} `json:"data"`
	}

	_ = json.Unmarshal(respBody, &result)

	_ = d.Set("message", result.Data.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("error_code", result.Data.Error)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))
	d.SetId(fmt.Sprintf("%s-%s", groupID, time.Now().UTC().Format("20060102T150405Z")))

	return diags
}

func resourceAgentRestartGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceAgentRestartGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
