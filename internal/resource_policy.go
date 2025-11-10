package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// wazuh_policy <-> /security/policies
//   - POST /security/policies
//   - GET  /security/policies
//   - PUT  /security/policies/{policy_id}
//   - DELETE /security/policies
func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},

		Schema: map[string]*schema.Schema{
			// User-facing fields
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name (≤ 64 characters). Should be unique to allow reliable lookup on creation.",
			},
			// JSON-encoded policy definition
			"policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON-encoded policy definition as expected by Wazuh (e.g. via jsonencode()).",
			},

			// Computed Wazuh ID
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric Wazuh policy ID as returned by the Wazuh API.",
			},
		},
	}
}

// helper struct for list responses
type wazuhPolicyListResponse struct {
	Data struct {
		AffectedItems []struct {
			ID     int             `json:"id"`
			Name   string          `json:"name"`
			Policy json.RawMessage `json:"policy"`
		} `json:"affected_items"`
		TotalAffected int `json:"total_affected_items"`
	} `json:"data"`
	Message string `json:"message"`
	Error   int    `json:"error"`
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	name := strings.TrimSpace(d.Get("name").(string))
	if name == "" {
		return diag.Errorf("name must not be empty")
	}

	policyStr := strings.TrimSpace(d.Get("policy").(string))
	if policyStr == "" {
		return diag.Errorf("policy must not be empty (JSON string expected)")
	}

	var policyBody json.RawMessage
	if err := json.Unmarshal([]byte(policyStr), &policyBody); err != nil {
		return diag.Errorf("invalid JSON in policy: %v", err)
	}

	payload := map[string]interface{}{
		"name":   name,
		"policy": json.RawMessage(policyBody),
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/policies", client.Endpoint)

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
		return diag.Errorf("failed to create policy '%s': status %d, body: %s", name, resp.StatusCode, string(respBody))
	}

	policyID, readDiags := lookupPolicyIDByName(ctx, client, name)
	if readDiags.HasError() {
		return readDiags
	}
	if policyID == "" {
		return diag.Errorf("policy '%s' created but not found during lookup", name)
	}

	d.SetId(policyID)
	_ = d.Set("policy_id", policyID)

	return resourcePolicyRead(ctx, d, meta)
}

// Read: GET /security/policies?policy_ids=<id>&limit=1
func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/policies", client.Endpoint)

	q := url.Values{}
	q.Set("policy_ids", id)
	q.Set("limit", "1")
	urlStr = urlStr + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
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
		return diag.Errorf("failed to read policy '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	var result wazuhPolicyListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse policy read response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		d.SetId("")
		return diags
	}

	item := result.Data.AffectedItems[0]
	policyID := strconv.Itoa(item.ID)

	_ = d.Set("policy_id", policyID)
	_ = d.Set("name", item.Name)

	if len(item.Policy) > 0 {
		var tmp interface{}
		if err := json.Unmarshal(item.Policy, &tmp); err == nil {
			if normalized, err := json.Marshal(tmp); err == nil {
				_ = d.Set("policy", string(normalized))
			} else {
				_ = d.Set("policy", string(item.Policy))
			}
		} else {
			_ = d.Set("policy", string(item.Policy))
		}
	}

	return diags
}

// Update: PUT /security/policies/{policy_id}
func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	policyID := d.Id()
	if policyID == "" {
		return diag.Errorf("cannot update policy without ID")
	}

	payload := make(map[string]interface{})

	if d.HasChange("name") {
		name := strings.TrimSpace(d.Get("name").(string))
		if name == "" {
			return diag.Errorf("name must not be empty")
		}
		payload["name"] = name
	}

	if d.HasChange("policy") {
		policyStr := strings.TrimSpace(d.Get("policy").(string))
		if policyStr == "" {
			return diag.Errorf("policy must not be empty when changed")
		}

		var policyBody json.RawMessage
		if err := json.Unmarshal([]byte(policyStr), &policyBody); err != nil {
			return diag.Errorf("invalid JSON in policy: %v", err)
		}
		payload["policy"] = json.RawMessage(policyBody)
	}

	if len(payload) == 0 {
		return resourcePolicyRead(ctx, d, meta)
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/policies/%s", client.Endpoint, policyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to update policy '%s': status %d, body: %s", policyID, resp.StatusCode, string(respBody))
	}

	return resourcePolicyRead(ctx, d, meta)
}

// Delete: DELETE /security/policies?policy_ids=<id>
func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	id := d.Id()
	if id == "" {
		return diags
	}

	urlStr := fmt.Sprintf("%s/security/policies", client.Endpoint)

	q := url.Values{}
	q.Set("policy_ids", id)
	urlStr = urlStr + "?" + q.Encode()

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to delete policy '%s': status %d, body: %s", id, resp.StatusCode, string(body))
	}

	d.SetId("")
	return diags
}

// Import: terraform import wazuh_policy.example 3
func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("policy_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func lookupPolicyIDByName(ctx context.Context, client *APIClient, name string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/security/policies", client.Endpoint)
	q := url.Values{}
	q.Set("search", name)
	q.Set("limit", "100")
	urlStr = urlStr + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return "", diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", diag.FromErr(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", diag.FromErr(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", diag.Errorf("failed to lookup policy '%s' after create: status %d, body: %s", name, resp.StatusCode, string(body))
	}

	var result wazuhPolicyListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", diag.Errorf("failed to parse policy lookup response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 || len(result.Data.AffectedItems) == 0 {
		return "", diag.Errorf("policy '%s' not found in list after create", name)
	}

	var matches []int
	for _, item := range result.Data.AffectedItems {
		if item.Name == name {
			matches = append(matches, item.ID)
		}
	}

	if len(matches) == 0 {
		return "", diag.Errorf("no exact policy name match found for '%s' after create", name)
	}
	if len(matches) > 1 {
		return "", diag.Errorf("multiple policies with name '%s' found; please ensure unique policy names", name)
	}

	return strconv.Itoa(matches[0]), diags
}
