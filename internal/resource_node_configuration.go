package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceNodeConfiguration manages Wazuh node configuration via /cluster/{node_id}/configuration.
func resourceNodeConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeConfigurationCreate,
		ReadContext:   resourceNodeConfigurationRead,
		UpdateContext: resourceNodeConfigurationUpdate,
		DeleteContext: resourceNodeConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNodeConfigurationImport,
		},

		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster node name (node_id) whose configuration is managed.",
			},
			"configuration_xml": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full XML content of the node ossec.conf.",
			},
		},
	}
}

// Create is effectively the same as Update: PUT /cluster/{node_id}/configuration
func resourceNodeConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNodeConfigurationPut(ctx, d, meta)
}

// Update node configuration (PUT)
func resourceNodeConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNodeConfigurationPut(ctx, d, meta)
}

// Shared PUT logic
func resourceNodeConfigurationPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	nodeID := d.Get("node_id").(string)
	xmlData := d.Get("configuration_xml").(string)

	escapedNode := url.PathEscape(nodeID)
	urlStr := fmt.Sprintf("%s/cluster/%s/configuration", client.Endpoint, escapedNode)

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer([]byte(xmlData)))
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

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to update configuration for node '%s': status %d, body: %s", nodeID, resp.StatusCode, string(body))
	}

	d.SetId(nodeID)
	return diags
}

// Read node configuration (GET /cluster/{node_id}/configuration?raw=true)
func resourceNodeConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	nodeID := d.Id()
	escapedNode := url.PathEscape(nodeID)
	urlStr := fmt.Sprintf("%s/cluster/%s/configuration?raw=true", client.Endpoint, escapedNode)

	req, err := http.NewRequest("GET", urlStr, nil)
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
		// Node configuration not found – remove from state
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read configuration for node '%s': status %d, body: %s", nodeID, resp.StatusCode, string(body))
	}

	// With raw=true we expect plain XML; store it back into state
	_ = d.Set("configuration_xml", string(body))
	_ = d.Set("node_id", nodeID)

	return diags
}

// Delete node configuration – not supported by Wazuh API, so this is a no-op (state only).
func resourceNodeConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

// Import existing node configuration by node_id
func resourceNodeConfigurationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	nodeID := d.Id()

	if err := d.Set("node_id", nodeID); err != nil {
		return nil, fmt.Errorf("failed to set node_id during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
