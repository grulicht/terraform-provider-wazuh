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

// resourceAgentUpgrade models an action that upgrades agents using a WPK file
// from the online repository via PUT /agents/upgrade.
func resourceAgentUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentUpgradeCreate,
		ReadContext:   resourceAgentUpgradeRead,
		// No Update – it's a one-shot action
		DeleteContext: resourceAgentUpgradeDelete,

		Schema: map[string]*schema.Schema{
			// Required according to the API: agents_list (comma-separated list or "all")
			"agents_list": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of agent IDs to upgrade (e.g. [\"001\", \"002\"]) or the keyword \"all\" to select all agents.",
			},

			// Optional upgrade options
			"wpk_repo": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "WPK repository URL/path.",
			},
			"upgrade_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Wazuh version to upgrade agents to (e.g. \"4.14.0\").",
			},
			"use_http": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Use HTTP instead of HTTPS when downloading the WPK. Default is false (HTTPS).",
			},
			"force": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Force upgrade, even if the agent appears to be already on the requested version.",
			},
			"package_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Package type to use for the upgrade (e.g. \"rpm\" or \"deb\"). By default, the manager infers this.",
			},

			// Computed API response
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after creating the upgrade tasks.",
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
				Description: "UTC timestamp when the upgrade request was sent.",
			},
		},
	}
}

// Create: send upgrade request via PUT /agents/upgrade
func resourceAgentUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/agents/upgrade", client.Endpoint)

	query := url.Values{}

	// agents_list is required (list of IDs or keyword "all")
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

	if v, ok := d.GetOk("wpk_repo"); ok {
		query.Set("wpk_repo", v.(string))
	}
	if v, ok := d.GetOk("upgrade_version"); ok {
		query.Set("upgrade_version", v.(string))
	}
	if v, ok := d.GetOk("package_type"); ok {
		query.Set("package_type", v.(string))
	}
	if v, ok := d.GetOkExists("use_http"); ok {
		if v.(bool) {
			query.Set("use_http", "true")
		} else {
			query.Set("use_http", "false")
		}
	}
	if v, ok := d.GetOkExists("force"); ok {
		if v.(bool) {
			query.Set("force", "true")
		} else {
			query.Set("force", "false")
		}
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	// Some endpoints accept empty body; we send {} for consistency
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
		return diag.Errorf("failed to upgrade agents: status %d, body: %s", resp.StatusCode, string(respBody))
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

	// Flatten affected_items into []map[string]interface{}
	ai := make([]map[string]interface{}, 0, len(result.Data.AffectedItems))
	for _, item := range result.Data.AffectedItems {
		ai = append(ai, map[string]interface{}{
			"agent":   item.Agent,
			"task_id": item.TaskID,
		})
	}
	_ = d.Set("affected_items", ai)

	t := time.Now().UTC()
	_ = d.Set("timestamp", t.Format(time.RFC3339))

	// Unique ID: timestamp-based
	d.SetId(t.Format("20060102T150405Z"))

	return diags
}

// Read: no-op – one-shot action
func resourceAgentUpgradeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: state-only, no API call
func resourceAgentUpgradeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
