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

// resourceCDBList defines the Terraform resource schema and CRUD operations
// for managing Wazuh CDB list files via /lists/files/{filename}.
func resourceCDBList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCDBListCreate,
		ReadContext:   resourceCDBListRead,
		UpdateContext: resourceCDBListUpdate,
		DeleteContext: resourceCDBListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCDBListImport,
		},

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CDB list filename (e.g. 'lists/mylist.cdb'). The Wazuh API searches this recursively.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CDB list file content to upload (plain text).",
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to overwrite the file if it already exists. Maps to the 'overwrite' query parameter.",
			},
		},
	}
}

// Create CDB list file (PUT /lists/files/{filename}?overwrite=...)
func resourceCDBListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCDBListPut(ctx, d, meta)
}

// Update CDB list file (same as Create: PUT /lists/files/{filename}?overwrite=...)
func resourceCDBListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceCDBListPut(ctx, d, meta)
}

// Shared logic for PUT
func resourceCDBListPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Get("filename").(string)
	content := d.Get("content").(string)
	overwrite := d.Get("overwrite").(bool)

	escaped := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/lists/files/%s?overwrite=%t", client.Endpoint, escaped, overwrite)

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer([]byte(content)))
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

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to upload CDB list file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	d.SetId(filename)
	return diags
}

// Read CDB list file (GET /lists/files/{filename}?raw=true)
func resourceCDBListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	escaped := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/lists/files/%s?raw=true", client.Endpoint, escaped)

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
		// File no longer exists — remove from state
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read CDB list file '%s': status %d, body: %s", filename, resp.StatusCode, string(body))
	}

	// body is plain text content when raw=true
	_ = d.Set("content", string(body))
	_ = d.Set("filename", filename)

	return diags
}

// Delete CDB list file (DELETE /lists/files/{filename})
func resourceCDBListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	escaped := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/lists/files/%s", client.Endpoint, escaped)

	req, err := http.NewRequest("DELETE", urlStr, nil)
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
	if resp.StatusCode != http.StatusNotFound && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return diag.Errorf("failed to delete CDB list file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import existing CDB list file by filename
func resourceCDBListImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	filename := d.Id()

	if err := d.Set("filename", filename); err != nil {
		return nil, fmt.Errorf("failed to set filename during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
