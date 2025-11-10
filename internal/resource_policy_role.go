package internal

import (
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

// resourcePolicyRole manages the relation between a Wazuh role and one or more policies
// via:
//   - POST   /security/roles/{role_id}/policies   (Create – link policies to role)
//   - DELETE /security/roles/{role_id}/policies   (Delete – unlink policies from role)
func resourcePolicyRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyRoleCreate,
		ReadContext:   resourcePolicyRoleRead,
		DeleteContext: resourcePolicyRoleDelete,

		// one-shot mapping resource – all arguments ForceNew
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh role ID to which the policies will be assigned.",
			},
			"policy_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of Wazuh policy IDs to link to this role.",
			},
			"position": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional security position for policies within the role.",
			},

			// Computed fields from API response
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after linking/unlinking policies.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of affected items reported by the Wazuh API.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of failed items reported by the Wazuh API.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw error code from the Wazuh API response (0 = success).",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp (RFC3339) when this mapping was last applied.",
			},
		},
	}
}

// Create: POST /security/roles/{role_id}/policies?policy_ids=1,2,3[&position=N]
func resourcePolicyRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	roleID := strings.TrimSpace(d.Get("role_id").(string))
	if roleID == "" {
		return diag.Errorf("role_id must not be empty")
	}

	rawPolicies := d.Get("policy_ids").([]interface{})
	if len(rawPolicies) == 0 {
		return diag.Errorf("policy_ids must contain at least one policy ID")
	}

	policyIDs := make([]string, 0, len(rawPolicies))
	for _, p := range rawPolicies {
		switch v := p.(type) {
		case int:
			policyIDs = append(policyIDs, fmt.Sprintf("%d", v))
		case int64:
			policyIDs = append(policyIDs, fmt.Sprintf("%d", v))
		case string:
			if strings.TrimSpace(v) != "" {
				policyIDs = append(policyIDs, strings.TrimSpace(v))
			}
		default:
			return diag.Errorf("unsupported type in policy_ids: %T", p)
		}
	}

	if len(policyIDs) == 0 {
		return diag.Errorf("policy_ids must contain at least one valid policy ID")
	}

	urlStr := fmt.Sprintf("%s/security/roles/%s/policies", client.Endpoint, roleID)

	query := url.Values{}
	query.Set("policy_ids", strings.Join(policyIDs, ","))

	if v, ok := d.GetOk("position"); ok {
		query.Set("position", fmt.Sprintf("%d", v.(int)))
	}

	if enc := query.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, nil)
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
		return diag.Errorf("failed to link policies %v to role %s: status %d, body: %s",
			policyIDs, roleID, resp.StatusCode, string(respBody))
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

	// Synthetic ID: role-<roleID>-policies-<comma-separated>
	d.SetId(fmt.Sprintf("role-%s-policies-%s", roleID, strings.Join(policyIDs, ",")))

	return diags
}

// Read: no-op – we do not re-query Wazuh for this mapping
func resourcePolicyRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: DELETE /security/roles/{role_id}/policies?policy_ids=1,2,3
func resourcePolicyRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	roleID := strings.TrimSpace(d.Get("role_id").(string))
	if roleID == "" {
		// nothing to do, just drop from state
		d.SetId("")
		return diags
	}

	rawPolicies := d.Get("policy_ids").([]interface{})
	policyIDs := make([]string, 0, len(rawPolicies))
	for _, p := range rawPolicies {
		switch v := p.(type) {
		case int:
			policyIDs = append(policyIDs, fmt.Sprintf("%d", v))
		case int64:
			policyIDs = append(policyIDs, fmt.Sprintf("%d", v))
		case string:
			if strings.TrimSpace(v) != "" {
				policyIDs = append(policyIDs, strings.TrimSpace(v))
			}
		}
	}

	urlStr := fmt.Sprintf("%s/security/roles/%s/policies", client.Endpoint, roleID)

	query := url.Values{}
	if len(policyIDs) > 0 {
		query.Set("policy_ids", strings.Join(policyIDs, ","))
	} else {
		// for safety, if list is empty we *don't* call with `all`
		// and just forget the state
		d.SetId("")
		return diags
	}

	if enc := query.Encode(); enc != "" {
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
		return diag.Errorf("failed to unlink policies %v from role %s: status %d, body: %s",
			policyIDs, roleID, resp.StatusCode, string(body))
	}

	d.SetId("")
	return diags
}
