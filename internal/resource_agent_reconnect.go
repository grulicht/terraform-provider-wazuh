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

// resourceAgentReconnect models an action that forces agents to reconnect
// via PUT /agents/reconnect.
func resourceAgentReconnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentReconnectCreate,
		ReadContext:   resourceAgentReconnectRead,
		// No Update – it's a one-shot action
		DeleteContext: resourceAgentReconnectDelete,

		Schema: map[string]*schema.Schema{
			"agents_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Optional list of agent IDs to force reconnect (e.g. [\"001\", \"002\"]). If omitted, all agents are targeted.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after sending the force reconnect command.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents for which the force reconnect command was processed.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents where the force reconnect command failed.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw error code returned by the Wazuh API (0 = success).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the force reconnect request was sent.",
			},
		},
	}
}

// Create: send force reconnect command via PUT /agents/reconnect
func resourceAgentReconnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/agents/reconnect", client.Endpoint)

	// Build query params
	query := url.Values{}

	if v, ok := d.GetOk("agents_list"); ok {
		raw := v.([]interface{})
		if len(raw) > 0 {
			ids := make([]string, 0, len(raw))
			for _, r := range raw {
				if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
					ids = append(ids, strings.TrimSpace(s))
				}
			}
			if len(ids) > 0 {
				// Wazuh expects agents_list=001,002,003
				query.Set("agents_list", strings.Join(ids, ","))
			}
		}
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

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
		return diag.Errorf("failed to force reconnect agents: status %d, body: %s", resp.StatusCode, string(respBody))
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

	// Use timestamp as unique ID
	d.SetId(time.Now().UTC().Format("20060102T150405Z"))

	return diags
}

// Read: no-op – we don't re-query Wazuh for this one-shot action
func resourceAgentReconnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: just remove from state, no API call
func resourceAgentReconnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
