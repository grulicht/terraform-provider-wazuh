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

// resourceUser manages Wazuh API users via:
//   - POST   /security/users         (Create)
//   - GET    /security/users         (Read, list/filter by user_ids)
//   - PUT    /security/users/{id}    (Update password)
//   - DELETE /security/users         (Delete by user_ids)
func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			// terraform import wazuh_user.example <user_id>
			StateContext: resourceUserImport,
		},

		Schema: map[string]*schema.Schema{
			// ---- Inputs ----

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh API username (4–64 characters).",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the Wazuh API user. On update, changing this will rotate the user's password.",
			},

			// ---- Computed fields ----
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric Wazuh user ID. This matches the Terraform resource ID.",
			},
		},
	}
}

// ----------------------
// Create: POST /security/users
// ----------------------
func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	username := strings.TrimSpace(d.Get("username").(string))
	password := d.Get("password").(string)

	if username == "" {
		return diag.Errorf("username must not be empty")
	}
	if len(username) < 4 || len(username) > 64 {
		return diag.Errorf("username must be between 4 and 64 characters")
	}
	if password == "" {
		return diag.Errorf("password must be provided when creating a Wazuh user")
	}

	payload := map[string]string{
		"username": username,
		"password": password,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/users", client.Endpoint)

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
		return diag.Errorf("failed to create Wazuh user '%s': status %d, body: %s",
			username, resp.StatusCode, string(respBody))
	}

	userID, err := findUserIDByUsername(ctx, client, username)
	if err != nil {
		return diag.FromErr(err)
	}

	if userID == "" {
		return diag.Errorf("user '%s' created but user_id could not be determined", username)
	}

	d.SetId(userID)
	_ = d.Set("user_id", userID)

	return resourceUserRead(ctx, d, meta)
}

// ----------------------
// Read: GET /security/users?user_ids=<id>
// ----------------------
func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/users", client.Endpoint)

	q := url.Values{}
	q.Set("user_ids", id)
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
		return diag.Errorf("failed to read Wazuh user '%s': status %d, body: %s",
			id, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID       json.Number `json:"id"`
				Username string      `json:"username"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse Wazuh user read response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		d.SetId("")
		return diags
	}

	item := result.Data.AffectedItems[0]
	userID := item.ID.String()

	d.SetId(userID)
	_ = d.Set("user_id", userID)
	_ = d.Set("username", item.Username)

	return diags
}

// ----------------------
// Update: PUT /security/users/{user_id} (password change)
// ----------------------
func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	id := d.Id()
	if id == "" {
		return diag.Errorf("cannot update Wazuh user: missing ID")
	}

	if d.HasChange("password") {
		newPass := d.Get("password").(string)
		if newPass == "" {
			return diag.Errorf("password cannot be empty when updating Wazuh user")
		}

		payload := map[string]string{
			"password": newPass,
		}
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return diag.FromErr(err)
		}

		urlStr := fmt.Sprintf("%s/security/users/%s", client.Endpoint, id)

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
			return diag.Errorf("failed to update Wazuh user '%s' password: status %d, body: %s",
				id, resp.StatusCode, string(respBody))
		}
	}

	return resourceUserRead(ctx, d, meta)
}

// ----------------------
// Delete: DELETE /security/users?user_ids=<id>
// ----------------------
func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/users", client.Endpoint)

	q := url.Values{}
	q.Set("user_ids", id)
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
		return diag.Errorf("failed to delete Wazuh user '%s': status %d, body: %s",
			id, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// ----------------------
// Import: terraform import wazuh_user.example <user_id>
// ----------------------
func resourceUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// ID from import is the Wazuh numeric user_id
	_ = d.Set("user_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

// ----------------------
// Helper: find user_id by username via GET /security/users
// ----------------------
func findUserIDByUsername(ctx context.Context, client *APIClient, username string) (string, error) {
	urlStr := fmt.Sprintf("%s/security/users", client.Endpoint)

	q := url.Values{}
	q.Set("search", username)
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
		return "", fmt.Errorf("failed to lookup Wazuh user '%s': status %d, body: %s",
			username, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID       json.Number `json:"id"`
				Username string      `json:"username"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse user lookup response: %w", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 {
		return "", fmt.Errorf("no users found matching username '%s'", username)
	}

	var matchIDs []string
	for _, item := range result.Data.AffectedItems {
		if item.Username == username {
			matchIDs = append(matchIDs, item.ID.String())
		}
	}

	if len(matchIDs) == 0 {
		return "", fmt.Errorf("no exact user match for username '%s'", username)
	}
	if len(matchIDs) > 1 {
		return "", fmt.Errorf("multiple users found with username '%s'; cannot determine unique user_id", username)
	}

	return matchIDs[0], nil
}
