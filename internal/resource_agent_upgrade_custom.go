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

// resourceAgentUpgradeCustom models an action that upgrades agents using a
// local WPK file via PUT /agents/upgrade_custom.
func resourceAgentUpgradeCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentUpgradeCustomCreate,
		ReadContext:   resourceAgentUpgradeCustomRead,
		// One-shot action – no Update
		DeleteContext: resourceAgentUpgradeCustomDelete,

		Schema: map[string]*schema.Schema{
			// Required: agents_list (IDs or "all")
			"agents_list": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of agent IDs to upgrade (e.g. [\"001\", \"002\"]) or the keyword \"all\" to select all agents.",
			},

			// Required: file_path – full path to WPK file on the manager
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Full path to the WPK file on the Wazuh manager (must be under Wazuh installation directory, e.g. /var/ossec).",
			},

			// Optional: installer script
			"installer": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Installation script to use (e.g. upgrade.sh or upgrade.bat). If omitted, Wazuh defaults apply.",
			},

			// Computed API response
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after creating the custom upgrade tasks.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents for which upgrade tasks were created.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents where upgrade tasks could not be created.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw error code returned by Wazuh API (0 = success).",
			},
			"affected_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of affected agents and their corresponding upgrade task IDs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent ID.",
						},
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID of the created upgrade task.",
						},
					},
				},
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the custom upgrade request was sent.",
			},
		},
	}
}

// Create: send custom upgrade request via PUT /agents/upgrade_custom
func resourceAgentUpgradeCustomCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/agents/upgrade_custom", client.Endpoint)

	query := url.Values{}

	// agents_list is required
	rawAgents := d.Get("agents_list").([]interface{})
	if len(rawAgents) == 0 {
		return diag.Errorf("agents_list must contain at least one value (agent ID or \"all\")")
	}
	agents := make([]string, 0, len(rawAgents))
	for _, r := range rawAgents {
		if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
			agents = append(agents, strings.TrimSpace(s))
		}
	}
	if len(agents) == 0 {
		return diag.Errorf("agents_list must contain at least one non-empty value")
	}
	query.Set("agents_list", strings.Join(agents, ","))

	// file_path is required
	filePath := strings.TrimSpace(d.Get("file_path").(string))
	if filePath == "" {
		return diag.Errorf("file_path must be a non-empty string")
	}
	query.Set("file_path", filePath)

	// optional installer
	if v, ok := d.GetOk("installer"); ok {
		query.Set("installer", v.(string))
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	// Minimal JSON body
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
		return diag.Errorf("failed to perform custom upgrade: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				Agent  string `json:"agent"`
				TaskID int    `json:"task_id"`
			} `json:"affected_items"`
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

	ai := make([]map[string]interface{}, 0, len(result.Data.AffectedItems))
	for _, item := range result.Data.AffectedItems {
		ai = append(ai, map[string]interface{}{
			"agent":   item.Agent,
			"task_id": item.TaskID,
		})
	}
	_ = d.Set("affected_items", ai)

	now := time.Now().UTC()
	_ = d.Set("timestamp", now.Format(time.RFC3339))

	// ID = timestamp-based
	d.SetId(now.Format("20060102T150405Z"))

	return diags
}

// Read: no-op – action resource
func resourceAgentUpgradeCustomRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: state-only, no API call
func resourceAgentUpgradeCustomDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
