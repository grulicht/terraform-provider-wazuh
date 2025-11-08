package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceManagerConfiguration manages the Wazuh manager configuration (ossec.conf)
// via:
//   - GET /manager/configuration
//   - PUT /manager/configuration
func resourceManagerConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceManagerConfigurationCreate,
		ReadContext:   resourceManagerConfigurationRead,
		UpdateContext: resourceManagerConfigurationUpdate,
		DeleteContext: resourceManagerConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceManagerConfigurationImport,
		},

		Schema: map[string]*schema.Schema{
			"configuration_xml": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full ossec.conf XML content for the Wazuh manager configuration.",
			},
			"last_updated_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the manager configuration was last updated via Terraform.",
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after updating the manager configuration.",
			},
		},
	}
}

// Create: upload manager configuration via PUT /manager/configuration
func resourceManagerConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return upsertManagerConfiguration(ctx, d, meta)
}

// Update: same as Create – replaces the manager configuration
func resourceManagerConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return upsertManagerConfiguration(ctx, d, meta)
}

func upsertManagerConfiguration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	configXML := d.Get("configuration_xml").(string)

	urlStr := fmt.Sprintf("%s/manager/configuration", client.Endpoint)

	// You can optionally add pretty/wait_for_complete here if needed
	query := url.Values{}
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", urlStr, bytes.NewBuffer([]byte(configXML)))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)
	req.Header.Set("Content-Type", "application/octet-stream")

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
		return diag.Errorf("failed to update manager configuration: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	}

	_ = json.Unmarshal(respBody, &result)

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error while updating manager configuration: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("last_updated_timestamp", time.Now().UTC().Format(time.RFC3339))

	// Single manager config – use a fixed ID
	d.SetId("manager")

	return diags
}

// Read: get current manager configuration via GET /manager/configuration?raw=true
func resourceManagerConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/manager/configuration", client.Endpoint)

	query := url.Values{}
	// raw = plain ossec.conf instead of structured JSON
	query.Set("raw", "true")

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
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
		// No configuration? (unusual) – drop from state
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read manager configuration: status %d, body: %s", resp.StatusCode, string(body))
	}

	// raw=true -> plain ossec.conf content
	_ = d.Set("configuration_xml", string(body))

	// ensure ID is set
	if d.Id() == "" {
		d.SetId("manager")
	}

	return diags
}

// Delete: no API call (there's no "delete configuration" endpoint). Just remove from state.
func resourceManagerConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

// Import: attach the existing manager configuration into Terraform state.
// ID is ignored; we always treat it as "manager".
func resourceManagerConfigurationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// normalize ID to "manager"
	d.SetId("manager")
	return []*schema.ResourceData{d}, nil
}
