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

// resourceGroup defines the Terraform resource schema and CRUD operations
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Wazuh group name. Must match pattern [a-zA-Z0-9_.-] and not be '.' or '..'.",
			},
		},
	}
}

// Create a new Wazuh group
func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Get("group_id").(string)
	payload := map[string]string{"group_id": groupID}

	body, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("%s/groups", client.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to create group '%s': status %d, body: %s", groupID, resp.StatusCode, string(respBody))
	}

	d.SetId(groupID)
	return diags
}

// Read a Wazuh group (GET /groups?groups_list=<group_id>)
func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Id()
	url := fmt.Sprintf("%s/groups?groups_list=%s", client.Endpoint, groupID)

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
		d.SetId("") // group not found
		return diags
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read group '%s': status %d, body: %s", groupID, resp.StatusCode, string(body))
	}

	// Parse the response
	var result struct {
		Data struct {
			AffectedItems []map[string]interface{} `json:"affected_items"`
			TotalAffected int                      `json:"total_affected_items"`
		} `json:"data"`
		Error int `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return diag.Errorf("failed to parse response: %v", err)
	}

	if result.Error != 0 || result.Data.TotalAffected == 0 {
		// group not found
		d.SetId("")
		return diags
	}

	// Optional: sync group_id again
	_ = d.Set("group_id", groupID)

	return diags
}

// Delete a Wazuh group (DELETE /groups?groups_list=<group_id>)
func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	groupID := d.Id()
	url := fmt.Sprintf("%s/groups?groups_list=%s", client.Endpoint, groupID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to delete group '%s': status %d, body: %s", groupID, resp.StatusCode, string(respBody))
	}

	// Parse API response (optional check)
	var result struct {
		Data struct {
			Error int `json:"error"`
		} `json:"data"`
	}
	_ = json.Unmarshal(respBody, &result)
	if result.Data.Error != 0 {
		return diag.Errorf("Wazuh API returned error while deleting group '%s': %s", groupID, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import existing group
func resourceGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	groupID := d.Id()

	// Set the attribute manually for Terraform state
	if err := d.Set("group_id", groupID); err != nil {
		return nil, fmt.Errorf("failed to set group_id for import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
