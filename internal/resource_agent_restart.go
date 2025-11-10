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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAgentRestart models an action that restarts all agents or a list of agents
// via PUT /agents/restart or PUT /agents/{agent_id}/restart (if exactly one agent).
func resourceAgentRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentRestartCreate,
		ReadContext:   resourceAgentRestartRead,
		// No Update, it's a one-shot action
		DeleteContext: resourceAgentRestartDelete,

		Schema: map[string]*schema.Schema{
			"agents_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Optional list of agent IDs to restart (e.g. [\"001\", \"002\"]). If omitted, all agents are restarted.",
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
				Description: "Raw error code returned by Wazuh API (0 = success, 2 = partial failure, etc.).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the restart request was sent.",
			},
		},
	}
}

// Create: send restart request via PUT /agents/restart
// or PUT /agents/{agent_id}/restart if only one agent is provided.
func resourceAgentRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	baseURL := client.Endpoint

	var (
		urlStr    string
		ids       []string
		hasAgents bool
	)

	// agents_list (optional)
	if v, ok := d.GetOk("agents_list"); ok {
		raw := v.([]interface{})
		if len(raw) > 0 {
			tmp := make([]string, 0, len(raw))
			for _, r := range raw {
				if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
					tmp = append(tmp, strings.TrimSpace(s))
				}
			}
			if len(tmp) > 0 {
				ids = tmp
				hasAgents = true
			}
		}
	}

	if !hasAgents {
		urlStr = fmt.Sprintf("%s/agents/restart", baseURL)
	} else if len(ids) == 1 {
		urlStr = fmt.Sprintf("%s/agents/%s/restart", baseURL, ids[0])
	} else {
		urlStr = fmt.Sprintf("%s/agents/restart", baseURL)
		query := url.Values{}
		query.Set("agents_list", strings.Join(ids, ","))
		if q := query.Encode(); q != "" {
			urlStr = urlStr + "?" + q
		}
	}

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
		return diag.Errorf("failed to restart agents: status %d, body: %s", resp.StatusCode, string(respBody))
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

	// Use timestamp as ID to make each restart action unique
	d.SetId(time.Now().UTC().Format("20060102T150405Z"))

	return diags
}

// Read: no-op, we don't re-call the API
func resourceAgentRestartRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: no API call – just forget this restart action from state
func resourceAgentRestartDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
