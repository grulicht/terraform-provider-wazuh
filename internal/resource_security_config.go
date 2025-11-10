package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSecurityConfig manages Wazuh security configuration via:
//   - GET    /security/config
//   - PUT    /security/config
//   - DELETE /security/config
func resourceSecurityConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityConfigCreate,
		ReadContext:   resourceSecurityConfigRead,
		UpdateContext: resourceSecurityConfigUpdate,
		DeleteContext: resourceSecurityConfigDelete,

		Schema: map[string]*schema.Schema{
			"auth_token_exp_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Time in seconds until the JWT token expires (>= 30).",
			},
			"rbac_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "RBAC mode: 'white' or 'black'.",
			},
		},
	}
}

// helper: build payload for PUT /security/config
func buildSecurityConfigPayload(d *schema.ResourceData) (map[string]interface{}, bool) {
	payload := make(map[string]interface{})

	if v, ok := d.GetOk("auth_token_exp_timeout"); ok {
		timeout := v.(int)
		if timeout > 0 {
			payload["auth_token_exp_timeout"] = timeout
		}
	}

	if v, ok := d.GetOk("rbac_mode"); ok {
		mode := v.(string)
		if mode != "" {
			payload["rbac_mode"] = mode
		}
	}

	return payload, len(payload) > 0
}

func resourceSecurityConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	payload, hasFields := buildSecurityConfigPayload(d)

	if !hasFields {
		d.SetId("security_config")
		return resourceSecurityConfigRead(ctx, d, meta)
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/security/config", client.Endpoint)

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
		return diag.Errorf("failed to update security config: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	}
	_ = json.Unmarshal(respBody, &result)

	if result.Error != 0 {
		return diag.Errorf("security config returned error code %d (%s)", result.Error, result.Message)
	}

	// Singleton ID
	d.SetId("security_config")
	if v, ok := payload["auth_token_exp_timeout"]; ok {
		_ = d.Set("auth_token_exp_timeout", v)
	}
	if v, ok := payload["rbac_mode"]; ok {
		_ = d.Set("rbac_mode", v)
	}

	return nil
}

// Read: GET /security/config
func resourceSecurityConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/security/config", client.Endpoint)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read security config: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			AuthTokenExpTimeout int    `json:"auth_token_exp_timeout"`
			RBACMode            string `json:"rbac_mode"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse security config response: %v", err)
	}

	if result.Error != 0 {
		return diag.Errorf("security config returned error code %d", result.Error)
	}

	_ = d.Set("auth_token_exp_timeout", result.Data.AuthTokenExpTimeout)
	_ = d.Set("rbac_mode", result.Data.RBACMode)

	if d.Id() == "" {
		d.SetId("security_config")
	}

	return diags
}

// Update: PUT /security/config
func resourceSecurityConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceSecurityConfigCreate(ctx, d, meta)
}

// Delete: DELETE /security/config (restore defaults)
func resourceSecurityConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/security/config", client.Endpoint)

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
		return diag.Errorf("failed to restore default security config: status %d, body: %s", resp.StatusCode, string(body))
	}

	// { "message": "Configuration was successfully updated", "error": 0 }
	d.SetId("")
	return diags
}
