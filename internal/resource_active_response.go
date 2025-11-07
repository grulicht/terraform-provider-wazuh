package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceActiveResponse defines the Terraform resource schema and CRUD operations for Wazuh active response.
func resourceActiveResponse() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveResponseCreate,
		ReadContext:   resourceActiveResponseRead,
		UpdateContext: resourceActiveResponseNoop,
		DeleteContext: resourceActiveResponseNoop,

		Importer: &schema.ResourceImporter{
			StateContext: resourceActiveResponseImport,
		},

		Schema: map[string]*schema.Schema{
			"command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Active Response command to run (as defined in Wazuh). If it starts with '!', it refers to a script name.",
			},
			"arguments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of arguments for the Active Response command.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"agents_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of agent IDs on which to run the command. If empty, all agents are targeted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response message from Wazuh after the command execution.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of agents affected by the command.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the command was executed.",
			},
		},
	}
}

// Create (run) the Active Response command via PUT /active-response
func resourceActiveResponseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	command := d.Get("command").(string)
	args := expandStringList(d.Get("arguments").([]interface{}))
	agents := expandStringList(d.Get("agents_list").([]interface{}))

	payload := map[string]interface{}{
		"command":   command,
		"arguments": args,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	// build base URL
	url := fmt.Sprintf("%s/active-response", client.Endpoint)

	// add agents_list as query parameter (comma-separated)
	if len(agents) > 0 {
		url = fmt.Sprintf("%s?agents_list=%s", url, joinComma(agents))
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
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
		return diag.Errorf("failed to execute active response '%s': status %d, body: %s", command, resp.StatusCode, string(respBody))
	}

	var result struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
		Data    struct {
			TotalAffected int `json:"total_affected_items"`
		} `json:"data"`
	}
	_ = json.Unmarshal(respBody, &result)

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use timestamp as unique ID
	d.SetId(fmt.Sprintf("%s-%d", command, time.Now().Unix()))

	return diags
}

// joinComma joins string slice into comma-separated string
func joinComma(items []string) string {
	out := ""
	for i, v := range items {
		if i > 0 {
			out += ","
		}
		out += v
	}
	return out
}

// Read (no actual data to fetch, only retain computed fields)
func resourceActiveResponseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// No need to contact API — this resource represents an executed command
	return diags
}

// Noop functions for update/delete
func resourceActiveResponseNoop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Import (simply reattach existing ID)
func resourceActiveResponseImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

// Helper to convert []interface{} → []string
func expandStringList(list []interface{}) []string {
	out := make([]string, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}
