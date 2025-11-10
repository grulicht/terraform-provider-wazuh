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

// wazuh_role_user – manages user ↔ roles assignments via:
//   - POST   /security/users/{user_id}/roles
//   - DELETE /security/users/{user_id}/roles
func resourceRoleUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleUserCreate,
		ReadContext:   resourceRoleUserRead,
		UpdateContext: resourceRoleUserUpdate,
		DeleteContext: resourceRoleUserDelete,

		// Import by "user_id:role_id1,role_id2"
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoleUserImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh user ID to which roles are assigned.",
			},
			"role_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of Wazuh role IDs assigned to the user.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"position": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Security position for roles/policies (Wazuh 'position' parameter). Applies when assigning roles.",
			},

			// Computed output from the last API call
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after assigning/removing roles.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of items affected by the last assign/remove operation.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of items that failed in the last assign/remove operation.",
			},
			"error_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Raw Wazuh API 'error' code of the last operation.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp of the last assign/remove operation.",
			},
		},
	}
}

// ---- Helpers ----

func expandStringToList(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, strings.TrimSpace(s))
		}
	}
	return out
}

type wazuhRoleUserResponse struct {
	Data struct {
		TotalAffected int `json:"total_affected_items"`
		TotalFailed   int `json:"total_failed_items"`
	} `json:"data"`
	Message string `json:"message"`
	Error   int    `json:"error"`
}

func setRoleUserResponseFields(d *schema.ResourceData, result *wazuhRoleUserResponse) {
	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("error_code", result.Error)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))
}

// ---- Create ----
// POST /security/users/{user_id}/roles?role_ids=1,2&position=<position>
func resourceRoleUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	userID := strings.TrimSpace(d.Get("user_id").(string))
	if userID == "" {
		return diag.Errorf("user_id must not be empty")
	}

	roleIDs := expandStringToList(d.Get("role_ids"))
	if len(roleIDs) == 0 {
		return diag.Errorf("role_ids must not be empty")
	}

	position := 0
	if v, ok := d.GetOk("position"); ok {
		position = v.(int)
	}

	if diags := callRoleUserAssign(ctx, client, userID, roleIDs, position, d); diags.HasError() {
		return diags
	}

	// ID resource: user_id + "|" + role_ids joined
	d.SetId(fmt.Sprintf("%s|%s", userID, strings.Join(roleIDs, ",")))
	return nil
}

func callRoleUserAssign(ctx context.Context, client *APIClient, userID string, roleIDs []string, position int, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/security/users/%s/roles", client.Endpoint, userID)

	q := url.Values{}
	q.Set("role_ids", strings.Join(roleIDs, ","))
	if position >= 0 {
		q.Set("position", fmt.Sprintf("%d", position))
	}
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	// body can be empty JSON
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to add roles %v to user '%s': status %d, body: %s",
			roleIDs, userID, resp.StatusCode, string(respBody))
	}

	var result wazuhRoleUserResponse
	_ = json.Unmarshal(respBody, &result)
	setRoleUserResponseFields(d, &result)

	return diags
}

// ---- Read ----
// No dedicated GET endpoint for user-roles mapping, so this is a no-op.
// We keep configuration + last response in state.
func resourceRoleUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// ---- Update ----
// Diff role_ids:
//   - newly added => POST /security/users/{user_id}/roles
//   - removed     => DELETE /security/users/{user_id}/roles
func resourceRoleUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	userID := strings.TrimSpace(d.Get("user_id").(string))
	if userID == "" {
		return diag.Errorf("user_id must not be empty")
	}

	if d.HasChange("role_ids") {
		oldRaw, newRaw := d.GetChange("role_ids")
		oldList := expandStringToList(oldRaw)
		newList := expandStringToList(newRaw)

		oldSet := make(map[string]struct{}, len(oldList))
		for _, v := range oldList {
			oldSet[v] = struct{}{}
		}
		newSet := make(map[string]struct{}, len(newList))
		for _, v := range newList {
			newSet[v] = struct{}{}
		}

		var toAdd, toRemove []string
		for v := range newSet {
			if _, ok := oldSet[v]; !ok {
				toAdd = append(toAdd, v)
			}
		}
		for v := range oldSet {
			if _, ok := newSet[v]; !ok {
				toRemove = append(toRemove, v)
			}
		}

		position := 0
		if v, ok := d.GetOk("position"); ok {
			position = v.(int)
		}

		// Assign new roles
		if len(toAdd) > 0 {
			if diags := callRoleUserAssign(ctx, client, userID, toAdd, position, d); diags.HasError() {
				return diags
			}
		}

		// Remove roles
		if len(toRemove) > 0 {
			if diags := callRoleUserRemove(ctx, client, userID, toRemove, d); diags.HasError() {
				return diags
			}
		}

		// Update ID to reflect new full set of roles
		roleIDs := expandStringToList(d.Get("role_ids"))
		d.SetId(fmt.Sprintf("%s|%s", userID, strings.Join(roleIDs, ",")))
	}

	// position change alone does not re-apply unless role_ids changed
	return nil
}

// DELETE /security/users/{user_id}/roles?role_ids=1,2
func callRoleUserRemove(ctx context.Context, client *APIClient, userID string, roleIDs []string, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/security/users/%s/roles", client.Endpoint, userID)

	q := url.Values{}
	q.Set("role_ids", strings.Join(roleIDs, ","))
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, urlStr, nil)
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
		return diag.Errorf("failed to remove roles %v from user '%s': status %d, body: %s",
			roleIDs, userID, resp.StatusCode, string(respBody))
	}

	var result wazuhRoleUserResponse
	_ = json.Unmarshal(respBody, &result)
	setRoleUserResponseFields(d, &result)

	return diags
}

// ---- Delete ----
// On destroy we remove exactly the roles from current state.
func resourceRoleUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	userID := strings.TrimSpace(d.Get("user_id").(string))
	if userID == "" {
		d.SetId("")
		return diags
	}

	roleIDs := expandStringToList(d.Get("role_ids"))
	if len(roleIDs) == 0 {
		d.SetId("")
		return diags
	}

	if diags := callRoleUserRemove(ctx, client, userID, roleIDs, d); diags.HasError() {
		return diags
	}

	d.SetId("")
	return diags
}

// ---- Import ----
// terraform import wazuh_role_user.example "5:1,2,3"
//
//	user_id=5, role_ids=[1,2,3]
func resourceRoleUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	raw := d.Id()
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected import ID format %q, expected 'user_id:role_id1,role_id2,...'", raw)
	}

	userID := strings.TrimSpace(parts[0])
	rolesPart := strings.TrimSpace(parts[1])

	var roleIDs []string
	if rolesPart != "" {
		for _, r := range strings.Split(rolesPart, ",") {
			r = strings.TrimSpace(r)
			if r != "" {
				roleIDs = append(roleIDs, r)
			}
		}
	}

	_ = d.Set("user_id", userID)
	_ = d.Set("role_ids", roleIDs)

	// Normalize ID to same pattern as Create/Update
	d.SetId(fmt.Sprintf("%s|%s", userID, strings.Join(roleIDs, ",")))

	return []*schema.ResourceData{d}, nil
}
