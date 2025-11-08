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

// resourceDecoder defines the Terraform resource schema and CRUD operations
// for managing Wazuh decoder files via /decoders/files/{filename}.
func resourceDecoder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDecoderCreate,
		ReadContext:   resourceDecoderRead,
		UpdateContext: resourceDecoderUpdate,
		DeleteContext: resourceDecoderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDecoderImport,
		},

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Decoder XML filename (e.g. 'local_decoder.xml'). Only the filename, no absolute path.",
			},
			"relative_dirname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional relative directory name where the decoder resides (e.g. 'decoders/local'). Maps to 'relative_dirname' query parameter.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full XML content of the decoder file.",
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to overwrite the file if it already exists. Maps to 'overwrite' query parameter.",
			},
		},
	}
}

// Create decoder file (PUT /decoders/files/{filename}?overwrite=...)
func resourceDecoderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDecoderPut(ctx, d, meta)
}

// Update decoder file (same as Create: PUT)
func resourceDecoderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDecoderPut(ctx, d, meta)
}

// Shared PUT logic
func resourceDecoderPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Get("filename").(string)
	content := d.Get("content").(string)
	overwrite := d.Get("overwrite").(bool)
	relativeDir := d.Get("relative_dirname").(string)

	escapedFilename := url.PathEscape(filename)

	// build base URL
	urlStr := fmt.Sprintf("%s/decoders/files/%s", client.Endpoint, escapedFilename)

	// build query parameters
	query := url.Values{}
	query.Set("overwrite", fmt.Sprintf("%t", overwrite))
	if relativeDir != "" {
		query.Set("relative_dirname", relativeDir)
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

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
		return diag.Errorf("failed to upload decoder file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	// ID = filename (relative_dirname is a separate attribute)
	d.SetId(filename)
	return diags
}

// Read decoder file (GET /decoders/files/{filename}?raw=true&relative_dirname=...)
func resourceDecoderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	relativeDir := d.Get("relative_dirname").(string)

	escapedFilename := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/decoders/files/%s", client.Endpoint, escapedFilename)

	query := url.Values{}
	query.Set("raw", "true")
	if relativeDir != "" {
		query.Set("relative_dirname", relativeDir)
	}
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

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
		// Decoder no longer exists
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read decoder file '%s': status %d, body: %s", filename, resp.StatusCode, string(body))
	}

	// body is plain XML content when raw=true
	_ = d.Set("content", string(body))
	_ = d.Set("filename", filename)

	return diags
}

// Delete decoder file (DELETE /decoders/files/{filename}?relative_dirname=...)
func resourceDecoderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	relativeDir := d.Get("relative_dirname").(string)

	escapedFilename := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/decoders/files/%s", client.Endpoint, escapedFilename)

	query := url.Values{}
	if relativeDir != "" {
		query.Set("relative_dirname", relativeDir)
	}
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

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
		return diag.Errorf("failed to delete decoder file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import existing decoder file by filename
func resourceDecoderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	filename := d.Id()

	if err := d.Set("filename", filename); err != nil {
		return nil, fmt.Errorf("failed to set filename during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
