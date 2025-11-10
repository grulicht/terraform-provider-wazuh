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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceRole manages Wazuh roles via:
//   - POST   /security/roles         (Create)
//   - GET    /security/roles         (Read, filter by role_ids)
//   - PUT    /security/roles/{id}    (Update name)
//   - DELETE /security/roles         (Delete by role_ids)
func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Importer: &schema.ResourceImporter{
			// terraform import wazuh_role.example <role_id>
			StateContext: resourceRoleImport,
		},

		Schema: map[string]*schema.Schema{
			// ---- Inputs ----

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role name (max 64 characters).",
			},

			// ---- Computed ----

			"role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric Wazuh role ID. This matches the Terraform resource ID.",
			},
		},
	}
}

// ----------------------
// Create: POST /security/roles
// ----------------------
func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	name := strings.TrimSpace(d.Get("name").(string))
	if name == "" {
		return diag.Errorf("name must not be empty")
	}
	if len(name) > 64 {
		return diag.Errorf("role name must be <= 64 characters")
	}

	payload := map[string]string{
		"name": name,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/roles", client.Endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to create Wazuh role '%s': status %d, body: %s",
			name, resp.StatusCode, string(respBody))
	}

	roleID, err := findRoleIDByName(ctx, client, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if roleID == "" {
		return diag.Errorf("role '%s' created but role_id could not be determined", name)
	}

	d.SetId(roleID)
	_ = d.Set("role_id", roleID)

	return resourceRoleRead(ctx, d, meta)
}

// ----------------------
// Read: GET /security/roles?role_ids=<id>
// ----------------------
func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/roles", client.Endpoint)

	q := url.Values{}
	q.Set("role_ids", id)
	q.Set("limit", "1")
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
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

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return diags
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read Wazuh role '%s': status %d, body: %s",
			id, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID   json.Number `json:"id"`
				Name string      `json:"name"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse Wazuh role read response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		d.SetId("")
		return diags
	}

	item := result.Data.AffectedItems[0]
	roleID := item.ID.String()

	d.SetId(roleID)
	_ = d.Set("role_id", roleID)
	_ = d.Set("name", item.Name)

	return diags
}

// ----------------------
// Update: PUT /security/roles/{role_id}
// ----------------------
func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	id := d.Id()
	if id == "" {
		return diag.Errorf("cannot update Wazuh role: missing ID")
	}

	if d.HasChange("name") {
		newName := strings.TrimSpace(d.Get("name").(string))
		if newName == "" {
			return diag.Errorf("role name cannot be empty")
		}
		if len(newName) > 64 {
			return diag.Errorf("role name must be <= 64 characters")
		}

		payload := map[string]string{
			"name": newName,
		}

		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return diag.FromErr(err)
		}

		urlStr := fmt.Sprintf("%s/security/roles/%s", client.Endpoint, id)

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
			return diag.Errorf("failed to update Wazuh role '%s': status %d, body: %s",
				id, resp.StatusCode, string(respBody))
		}
	}

	return resourceRoleRead(ctx, d, meta)
}

// ----------------------
// Delete: DELETE /security/roles?role_ids=<id>
// ----------------------
func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/roles", client.Endpoint)

	q := url.Values{}
	q.Set("role_ids", id)
	if enc := q.Encode(); enc != "" {
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to delete Wazuh role '%s': status %d, body: %s",
			id, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// ----------------------
// Import: terraform import wazuh_role.example <role_id>
// ----------------------
func resourceRoleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("role_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

// ----------------------
// Helper: find role_id by name via GET /security/roles
// ----------------------
func findRoleIDByName(ctx context.Context, client *APIClient, name string) (string, error) {
	urlStr := fmt.Sprintf("%s/security/roles", client.Endpoint)

	q := url.Values{}
	q.Set("search", name)
	q.Set("limit", "100")
	if enc := q.Encode(); enc != "" {
		urlStr = urlStr + "?" + enc
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to lookup Wazuh role '%s': status %d, body: %s",
			name, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID   json.Number `json:"id"`
				Name string      `json:"name"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse role lookup response: %w", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 {
		return "", fmt.Errorf("no roles found matching name '%s'", name)
	}

	var matchIDs []string
	for _, item := range result.Data.AffectedItems {
		if item.Name == name {
			matchIDs = append(matchIDs, item.ID.String())
		}
	}

	if len(matchIDs) == 0 {
		return "", fmt.Errorf("no exact role match for name '%s'", name)
	}
	if len(matchIDs) > 1 {
		return "", fmt.Errorf("multiple roles found with name '%s'; cannot determine unique role_id", name)
	}

	return matchIDs[0], nil
}
