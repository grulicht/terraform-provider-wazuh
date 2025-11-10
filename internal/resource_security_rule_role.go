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

// resourceSecurityRuleRole manages the relation between a role and security rules:
//
//   - POST   /security/roles/{role_id}/rules   (Create)
//   - DELETE /security/roles/{role_id}/rules   (Delete)
//
// It's a one-shot "link/unlink" relation resource, similar to wazuh_policy_role.
func resourceSecurityRuleRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityRuleRoleCreate,
		ReadContext:   resourceSecurityRuleRoleRead,
		DeleteContext: resourceSecurityRuleRoleDelete,

		// Relationship-type resource → ForceNew on all inputs
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Wazuh role to which security rules will be linked.",
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of security rule IDs to link to this role.",
			},

			// Optional: store some metadata from the response
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of security rules for which the relation was processed.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of security rules where the relation failed.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw error code returned by the Wazuh API (0 = success).",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after linking/unlinking rules.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the relation was last applied.",
			},
		},
	}
}

// Create: POST /security/roles/{role_id}/rules?rule_ids=1,2,3
func resourceSecurityRuleRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	roleID := strings.TrimSpace(d.Get("role_id").(string))
	if roleID == "" {
		return diag.Errorf("role_id must not be empty")
	}

	rawRuleIDs := d.Get("rule_ids").([]interface{})
	if len(rawRuleIDs) == 0 {
		return diag.Errorf("rule_ids must not be empty")
	}

	ruleIDStrs := make([]string, 0, len(rawRuleIDs))
	for _, v := range rawRuleIDs {
		ruleIDStrs = append(ruleIDStrs, fmt.Sprintf("%v", v))
	}

	u := fmt.Sprintf("%s/security/roles/%s/rules", client.Endpoint, roleID)

	q := url.Values{}
	q.Set("rule_ids", strings.Join(ruleIDStrs, ","))
	if enc := q.Encode(); enc != "" {
		u = u + "?" + enc
	}

	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to link security rules to role '%s': status %d, body: %s", roleID, resp.StatusCode, string(respBody))
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

	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("message", result.Message)
	_ = d.Set("error_code", result.Error)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Make a deterministic ID from role + rule_ids
	id := fmt.Sprintf("%s:%s", roleID, strings.Join(ruleIDStrs, ","))
	d.SetId(id)

	return diags
}

// Read: no-op – relationship resource, no dedicated GET for this specific mapping
func resourceSecurityRuleRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: DELETE /security/roles/{role_id}/rules?rule_ids=1,2,3
func resourceSecurityRuleRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	roleID := strings.TrimSpace(d.Get("role_id").(string))
	if roleID == "" {
		// fallback: try to parse from ID "roleID:rule1,rule2"
		parts := strings.SplitN(id, ":", 2)
		if len(parts) > 0 {
			roleID = parts[0]
		}
	}

	rawRuleIDs := d.Get("rule_ids").([]interface{})
	if len(rawRuleIDs) == 0 {
		// If somehow empty, nothing to unlink
		d.SetId("")
		return diags
	}

	ruleIDStrs := make([]string, 0, len(rawRuleIDs))
	for _, v := range rawRuleIDs {
		ruleIDStrs = append(ruleIDStrs, fmt.Sprintf("%v", v))
	}

	u := fmt.Sprintf("%s/security/roles/%s/rules", client.Endpoint, roleID)

	q := url.Values{}
	q.Set("rule_ids", strings.Join(ruleIDStrs, ","))
	if enc := q.Encode(); enc != "" {
		u = u + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return diag.FromErr(err)
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
		return diag.Errorf("failed to unlink security rules from role '%s': status %d, body: %s", roleID, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}
