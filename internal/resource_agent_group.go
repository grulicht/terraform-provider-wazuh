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

// resourceAgentGroup manages assignment of Wazuh agents to groups.
//
// It supports two modes:
//
// 1) Per-agent mode (single agent):
//   - PUT    /agents/{agent_id}/group/{group_id}   (Create)
//   - DELETE /agents/{agent_id}/group/{group_id}   (Delete)
//
// 2) Bulk mode (multiple or all agents), when agent_id is NOT set:
//   - PUT    /agents/group?group_id=...&agents_list=...&force_single_group=...
//   - DELETE /agents/group?group_id=...&agents_list=...
//
// This resource models either:
//   - "agent X is a member of group Y"  (per-agent), or
//   - "a set of agents (or all) has been assigned/removed to/from group Y" (bulk).
func resourceAgentGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentGroupCreate,
		ReadContext:   resourceAgentGroupRead,
		DeleteContext: resourceAgentGroupDelete,

		Schema: map[string]*schema.Schema{
			// ---- Inputs ----

			// Single-agent mode: set this to target a specific agent.
			// Bulk mode: leave empty, use agents_list instead.
			"agent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Wazuh agent ID (e.g. \"001\"). If set, per-agent endpoints are used. If omitted, bulk /agents/group endpoints are used.",
			},

			// Bulk mode: list of agents to assign/remove to/from a group via /agents/group.
			// If omitted in Create, Wazuh API assigns ALL agents to the group.
			// In Delete (bulk), agents_list is required.
			"agents_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of agent IDs for bulk operations via /agents/group. If agent_id is set, this is ignored.",
			},

			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh group ID (group name) used in the assignment/removal operation.",
			},

			// For per-agent and bulk assign operations.
			"force_single_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, removes the agent(s) from all groups they belong to and assigns them only to the specified group.",
			},

			// ---- Computed outputs (from Wazuh API) ----

			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh about the assignment/removal operation.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of affected items according to the Wazuh API response.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of failed items according to the Wazuh API response.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Error code returned by Wazuh (0 = success, >0 indicates partial/failed states).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the last operation (assign/remove) was performed.",
			},
		},
	}
}

// Create:
//
// - If agent_id is set  => PUT /agents/{agent_id}/group/{group_id}?force_single_group=...
// - If agent_id is empty => PUT /agents/group?group_id=...&agents_list=...&force_single_group=...
func resourceAgentGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := strings.TrimSpace(d.Get("agent_id").(string))
	groupID := strings.TrimSpace(d.Get("group_id").(string))
	forceSingle := d.Get("force_single_group").(bool)

	if groupID == "" {
		return diag.Errorf("group_id must not be empty")
	}

	// Minimal body for endpoints that expect JSON
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	var urlStr string

	// Per-agent mode
	if agentID != "" {
		urlStr = fmt.Sprintf("%s/agents/%s/group/%s", client.Endpoint, agentID, groupID)
		if forceSingle {
			urlStr = urlStr + "?force_single_group=true"
		}
	} else {
		// Bulk mode: /agents/group
		urlStr = fmt.Sprintf("%s/agents/group", client.Endpoint)
		q := url.Values{}
		q.Set("group_id", groupID)

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
					q.Set("agents_list", strings.Join(ids, ","))
				}
			}
		}
		// If agents_list is not set in Create, Wazuh will assign ALL agents to the group.

		if forceSingle {
			q.Set("force_single_group", "true")
		}

		if enc := q.Encode(); enc != "" {
			urlStr = urlStr + "?" + enc
		}
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
		return diag.Errorf("failed to assign agents to group '%s': status %d, body: %s",
			groupID, resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			TotalAffected int    `json:"total_affected_items"`
			TotalFailed   int    `json:"total_failed_items"`
			Message       string `json:"message"`
			Error         int    `json:"error"`
		} `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}

	_ = json.Unmarshal(respBody, &result)

	msg := result.Data.Message
	if msg == "" {
		msg = result.Message
	}

	errCode := result.Data.Error
	if errCode == 0 && result.Error != 0 {
		errCode = result.Error
	}

	_ = d.Set("message", msg)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("error_code", errCode)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	if agentID != "" {
		d.SetId(fmt.Sprintf("%s-%s", agentID, groupID))
	} else {
		d.SetId(fmt.Sprintf("%s-%s", groupID, time.Now().UTC().Format("20060102T150405Z")))
	}

	return diags
}

func resourceAgentGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceAgentGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	agentID := strings.TrimSpace(d.Get("agent_id").(string))
	groupID := strings.TrimSpace(d.Get("group_id").(string))

	if groupID == "" {
		d.SetId("")
		return diags
	}

	var urlStr string
	var req *http.Request
	var err error

	if agentID != "" {
		urlStr = fmt.Sprintf("%s/agents/%s/group/%s", client.Endpoint, agentID, groupID)
		req, err = http.NewRequestWithContext(ctx, "DELETE", urlStr, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		urlStr = fmt.Sprintf("%s/agents/group", client.Endpoint)
		q := url.Values{}
		q.Set("group_id", groupID)

		v, ok := d.GetOk("agents_list")
		if !ok {
			return diag.Errorf("agents_list must be provided for bulk delete when agent_id is not set")
		}
		raw := v.([]interface{})
		if len(raw) == 0 {
			return diag.Errorf("agents_list must not be empty for bulk delete when agent_id is not set")
		}

		ids := make([]string, 0, len(raw))
		for _, r := range raw {
			if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
				ids = append(ids, strings.TrimSpace(s))
			}
		}
		if len(ids) == 0 {
			return diag.Errorf("agents_list must contain at least one non-empty agent ID for bulk delete")
		}
		q.Set("agents_list", strings.Join(ids, ","))

		if enc := q.Encode(); enc != "" {
			urlStr = urlStr + "?" + enc
		}

		req, err = http.NewRequestWithContext(ctx, "DELETE", urlStr, nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

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
		if agentID != "" {
			return diag.Errorf("failed to remove agent '%s' from group '%s': status %d, body: %s",
				agentID, groupID, resp.StatusCode, string(respBody))
		}
		return diag.Errorf("failed to remove agents from group '%s': status %d, body: %s",
			groupID, resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			TotalAffected int    `json:"total_affected_items"`
			TotalFailed   int    `json:"total_failed_items"`
			Message       string `json:"message"`
			Error         int    `json:"error"`
		} `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}
	_ = json.Unmarshal(respBody, &result)

	msg := result.Data.Message
	if msg == "" {
		msg = result.Message
	}

	errCode := result.Data.Error
	if errCode == 0 && result.Error != 0 {
		errCode = result.Error
	}

	_ = d.Set("message", msg)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("error_code", errCode)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	d.SetId("")
	return diags
}
