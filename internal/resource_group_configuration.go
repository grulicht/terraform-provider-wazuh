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

// resourceGroupConfiguration defines the Terraform resource schema and CRUD operations for Wazuh group configuration.
func resourceGroupConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupConfigurationCreate,
		ReadContext:   resourceGroupConfigurationRead,
		UpdateContext: resourceGroupConfigurationUpdate,
		DeleteContext: resourceGroupConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupConfigurationImport,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh group ID (group name).",
			},
			"configuration_xml": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full XML configuration for the group (content of agent.conf).",
			},
		},
	}
}

// Create or update a Wazuh group configuration (PUT /groups/{group_id}/configuration)
func resourceGroupConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGroupConfigurationUpdate(ctx, d, meta)
}

// Update group configuration via PUT /groups/{group_id}/configuration
func resourceGroupConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Get("group_id").(string)
	xmlData := d.Get("configuration_xml").(string)

	url := fmt.Sprintf("%s/groups/%s/configuration", client.Endpoint, groupID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(xmlData)))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to update configuration for group '%s': status %d, body: %s", groupID, resp.StatusCode, string(body))
	}

	var result struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	}
	_ = json.Unmarshal(body, &result)
	if result.Error != 0 {
		return diag.Errorf("API error updating group configuration '%s': %s", groupID, string(body))
	}

	d.SetId(groupID)
	return diags
}

// Read group configuration (GET /groups/{group_id}/configuration)
func resourceGroupConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Id()
	url := fmt.Sprintf("%s/groups/%s/configuration", client.Endpoint, groupID)

	req, err := http.NewRequest("GET", url, nil)
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

	if resp.StatusCode == 404 {
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read group configuration '%s': status %d, body: %s", groupID, resp.StatusCode, string(body))
	}

	// Optional: Store the XML configuration back into state
	_ = d.Set("configuration_xml", string(body))
	_ = d.Set("group_id", groupID)

	return diags
}

// Delete group configuration — not supported (noop)
func resourceGroupConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Deleting a configuration is not supported by Wazuh API
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

// Import existing configuration
func resourceGroupConfigurationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	groupID := d.Id()

	if err := d.Set("group_id", groupID); err != nil {
		return nil, fmt.Errorf("failed to set group_id during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
