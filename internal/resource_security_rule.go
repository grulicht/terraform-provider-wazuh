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

// resourceSecurityRule manages Wazuh security rules via:
//   - POST   /security/rules        (Create)
//   - GET    /security/rules        (Read by rule_ids / search)
//   - PUT    /security/rules/{id}   (Update)
//   - DELETE /security/rules        (Delete by rule_ids)
func resourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityRuleCreate,
		ReadContext:   resourceSecurityRuleRead,
		UpdateContext: resourceSecurityRuleUpdate,
		DeleteContext: resourceSecurityRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityRuleImport,
		},

		Schema: map[string]*schema.Schema{
			// ---- Inputs ----

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security rule name (<= 64 characters).",
			},

			// We store the rule body as a JSON string in Terraform.
			// Example:
			// rule = jsonencode({
			//   MATCH = {
			//     definition = "normalRule"
			//   }
			// })
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON-encoded security rule body as expected by the Wazuh API.",
			},

			// ---- Computed ----

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric Wazuh security rule ID.",
			},
		},
	}
}

// Create: POST /security/rules
func resourceSecurityRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	name := strings.TrimSpace(d.Get("name").(string))
	if name == "" {
		return diag.Errorf("name must not be empty")
	}

	ruleStr := strings.TrimSpace(d.Get("rule").(string))
	if ruleStr == "" {
		return diag.Errorf("rule must not be empty (JSON string expected)")
	}

	// Parse rule JSON string into a generic object for the API body
	var ruleBody interface{}
	if err := json.Unmarshal([]byte(ruleStr), &ruleBody); err != nil {
		return diag.Errorf("invalid rule JSON: %v", err)
	}

	payload := map[string]interface{}{
		"name": name,
		"rule": ruleBody,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/rules", client.Endpoint)

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
		return diag.Errorf("failed to create security rule '%s': status %d, body: %s", name, resp.StatusCode, string(respBody))
	}

	// The create endpoint does NOT return rule_id directly.
	// We need to look it up via GET /security/rules?search=<name>
	ruleID, err := lookupSecurityRuleIDByName(ctx, client, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ruleID)
	_ = d.Set("rule_id", ruleID)

	// Sync from API (name + rule) if possible
	return resourceSecurityRuleRead(ctx, d, meta)
}

// Read: GET /security/rules?rule_ids=<id>
func resourceSecurityRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/rules", client.Endpoint)

	q := url.Values{}
	q.Set("rule_ids", id)
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
		// Rule no longer exists
		d.SetId("")
		return diags
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read security rule '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID   int         `json:"id"`
				Name string      `json:"name"`
				Rule interface{} `json:"rule"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse security rule read response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		// Not found / no items
		d.SetId("")
		return diags
	}

	item := result.Data.AffectedItems[0]

	_ = d.Set("rule_id", fmt.Sprintf("%d", item.ID))
	_ = d.Set("name", item.Name)

	// Re-marshal the rule object back into JSON string for the Terraform state
	if item.Rule != nil {
		if rb, err := json.Marshal(item.Rule); err == nil {
			_ = d.Set("rule", string(rb))
		}
	}

	return diags
}

// Update: PUT /security/rules/{rule_id}
func resourceSecurityRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	id := d.Id()
	if id == "" {
		return diag.Errorf("cannot update security rule: missing ID")
	}

	name := strings.TrimSpace(d.Get("name").(string))
	ruleStr := strings.TrimSpace(d.Get("rule").(string))

	payload := map[string]interface{}{}

	if name != "" {
		payload["name"] = name
	}

	if ruleStr != "" {
		var ruleBody interface{}
		if err := json.Unmarshal([]byte(ruleStr), &ruleBody); err != nil {
			return diag.Errorf("invalid rule JSON during update: %v", err)
		}
		payload["rule"] = ruleBody
	}

	// At least one field must be present (API requirement)
	if len(payload) == 0 {
		// Nothing to update
		return resourceSecurityRuleRead(ctx, d, meta)
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/rules/%s", client.Endpoint, id)

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
		return diag.Errorf("failed to update security rule '%s': status %d, body: %s", id, resp.StatusCode, string(respBody))
	}

	// Refresh state from API
	return resourceSecurityRuleRead(ctx, d, meta)
}

// Delete: DELETE /security/rules?rule_ids=<id>
func resourceSecurityRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/rules", client.Endpoint)

	q := url.Values{}
	q.Set("rule_ids", id)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to delete security rule '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	d.SetId("")
	return diags
}

// Import: terraform import wazuh_security_rule.example 5
func resourceSecurityRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// The import ID is the numeric rule_id
	_ = d.Set("rule_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

// Helper: lookup rule_id by unique name via GET /security/rules?search=<name>
func lookupSecurityRuleIDByName(ctx context.Context, client *APIClient, name string) (string, error) {
	urlStr := fmt.Sprintf("%s/security/rules", client.Endpoint)

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
		return "", fmt.Errorf("failed to lookup security rule by name '%s': status %d, body: %s", name, resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AffectedItems []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"affected_items"`
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse lookup response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		return "", fmt.Errorf("no security rule found with name '%s'", name)
	}

	// Filter for exact name match
	matches := []int{}
	for _, item := range result.Data.AffectedItems {
		if item.Name == name {
			matches = append(matches, item.ID)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no exact security rule name match found for '%s'", name)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple security rules found with name '%s', cannot determine unique rule_id", name)
	}

	return fmt.Sprintf("%d", matches[0]), nil
}
