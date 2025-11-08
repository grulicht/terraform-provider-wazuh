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

// resourceLogtest defines the Terraform resource schema and CRUD operations
// for running Wazuh logtest via PUT /logtest and ending the session via
// DELETE /logtest/sessions/{token}.
func resourceLogtest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogtestCreate,
		ReadContext:   resourceLogtestRead,
		UpdateContext: resourceLogtestNoop,
		DeleteContext: resourceLogtestDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLogtestImport,
		},

		Schema: map[string]*schema.Schema{
			"log_format": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log format for logtest (e.g. syslog, json, eventchannel, command, etc.).",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location string (path) used by logtest.",
			},
			"event": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Event/log line to evaluate.",
			},
			// Optional token: if provided, Wazuh may reuse an existing logtest session.
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Logtest session token. If omitted, a new session is created by Wazuh.",
			},
			"alert": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the event raised an alert.",
			},
			"codemsg": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Numeric code returned by logtest (e.g. 1 for no alert, 0 for success, etc.).",
			},
			"messages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of diagnostic messages returned by logtest.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"output": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Raw JSON output object from logtest, serialized as a string.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when logtest was executed.",
			},
		},
	}
}

// Create runs logtest via PUT /logtest
func resourceLogtestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	logFormat := d.Get("log_format").(string)
	location := d.Get("location").(string)
	event := d.Get("event").(string)

	payload := map[string]interface{}{
		"log_format": logFormat,
		"location":   location,
		"event":      event,
	}

	// Optional token reuse
	if v, ok := d.GetOk("token"); ok {
		if tokenStr := v.(string); tokenStr != "" {
			payload["token"] = tokenStr
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	urlStr := fmt.Sprintf("%s/logtest", client.Endpoint)
	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(body))
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
		return diag.Errorf("failed to run logtest: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Error int `json:"error"`
		Data  struct {
			Messages []string               `json:"messages"`
			Token    string                 `json:"token"`
			Output   map[string]interface{} `json:"output"`
			Alert    bool                   `json:"alert"`
			Codemsg  int                    `json:"codemsg"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return diag.Errorf("failed to parse logtest response: %v (body: %s)", err, string(respBody))
	}

	if result.Error != 0 {
		return diag.Errorf("Wazuh logtest returned error: %s", string(respBody))
	}

	// Serialize output object as JSON string for storage
	var outputJSON string
	if result.Data.Output != nil {
		if b, err := json.Marshal(result.Data.Output); err == nil {
			outputJSON = string(b)
		}
	}

	_ = d.Set("token", result.Data.Token)
	_ = d.Set("messages", result.Data.Messages)
	_ = d.Set("alert", result.Data.Alert)
	_ = d.Set("codemsg", result.Data.Codemsg)
	_ = d.Set("output", outputJSON)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use token as ID if present; otherwise fallback to timestamp-based ID
	if result.Data.Token != "" {
		d.SetId(result.Data.Token)
	} else {
		d.SetId(fmt.Sprintf("logtest-%d", time.Now().Unix()))
	}

	return diags
}

// Read: no-op. We keep the original logtest result in state.
func resourceLogtestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// No-op for update (re-running logtest would be a new action).
func resourceLogtestNoop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Delete: end the logtest session via DELETE /logtest/sessions/{token}
func resourceLogtestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	token := d.Id()
	if token == "" {
		// nothing to delete on server
		d.SetId("")
		return diags
	}

	escapedToken := url.PathEscape(token)
	urlStr := fmt.Sprintf("%s/logtest/sessions/%s", client.Endpoint, escapedToken)

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

	// 2xx = OK; 404 = already gone -> treat as success
	if resp.StatusCode != http.StatusNotFound && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return diag.Errorf("failed to delete logtest session '%s': status %d, body: %s", token, resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return diags
}

// Import: just attach an existing token as ID; no API call.
func resourceLogtestImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// The ID is assumed to be the token
	if err := d.Set("token", d.Id()); err != nil {
		return nil, fmt.Errorf("failed to set token during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
