package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceRule manages a Wazuh ruleset file via:
//   - PUT    /rules/files/{filename}
//   - GET    /rules/files/{filename}
//   - DELETE /rules/files/{filename}
func resourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the rules file (e.g. 'local_rules.xml').",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full XML content of the Wazuh rules file.",
			},
			"relative_dirname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional relative directory under the Wazuh rules path (e.g. 'rules').",
			},
			"overwrite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to overwrite an existing rules file with the same name.",
			},
		},
	}
}

// Create: upload/replace a rules file via PUT /rules/files/{filename}
func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return upsertRule(ctx, d, meta)
}

// Update: same behavior as Create – overwrite the rules file if allowed
func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return upsertRule(ctx, d, meta)
}

func upsertRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Get("filename").(string)
	content := d.Get("content").(string)
	overwrite := d.Get("overwrite").(bool)

	escapedFilename := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/rules/files/%s", client.Endpoint, escapedFilename)

	query := url.Values{}
	// overwrite query param
	query.Set("overwrite", strconv.FormatBool(overwrite))

	// optional relative_dirname
	if v, ok := d.GetOk("relative_dirname"); ok {
		rel := v.(string)
		if rel != "" {
			query.Set("relative_dirname", rel)
		}
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", urlStr, bytes.NewBuffer([]byte(content)))
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
		return diag.Errorf("failed to upload rules file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	// Use filename as the resource ID
	d.SetId(filename)

	return diags
}

// Read: fetch rules file content via GET /rules/files/{filename}?raw=true
func resourceRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	escapedFilename := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/rules/files/%s", client.Endpoint, escapedFilename)

	query := url.Values{}
	query.Set("raw", "true")

	if v, ok := d.GetOk("relative_dirname"); ok {
		rel := v.(string)
		if rel != "" {
			query.Set("relative_dirname", rel)
		}
	}

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
		// Rules file no longer exists
		d.SetId("")
		return diags
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("failed to read rules file '%s': status %d, body: %s", filename, resp.StatusCode, string(body))
	}

	// With raw=true, response body should be the plain XML content
	_ = d.Set("content", string(body))
	_ = d.Set("filename", filename)

	return diags
}

// Delete: delete rules file via DELETE /rules/files/{filename}
func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	filename := d.Id()
	escapedFilename := url.PathEscape(filename)
	urlStr := fmt.Sprintf("%s/rules/files/%s", client.Endpoint, escapedFilename)

	query := url.Values{}
	if v, ok := d.GetOk("relative_dirname"); ok {
		rel := v.(string)
		if rel != "" {
			query.Set("relative_dirname", rel)
		}
	}

	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
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

	respBody, _ := io.ReadAll(resp.Body)

	// 404 = already deleted → treat as success
	if resp.StatusCode != http.StatusNotFound && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return diag.Errorf("failed to delete rules file '%s': status %d, body: %s", filename, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import: attach existing rules file by filename (no API call here)
func resourceRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	filename := d.Id()
	if err := d.Set("filename", filename); err != nil {
		return nil, fmt.Errorf("failed to set filename during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
