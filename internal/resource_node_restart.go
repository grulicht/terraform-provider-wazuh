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

// resourceNodeRestart defines the Terraform resource schema and operations
// for restarting Wazuh cluster nodes via PUT /cluster/restart.
func resourceNodeRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeRestartCreate,
		ReadContext:   resourceNodeRestartRead,
		UpdateContext: resourceNodeRestartNoop,
		DeleteContext: resourceNodeRestartNoop,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNodeRestartImport,
		},

		Schema: map[string]*schema.Schema{
			"nodes_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of node IDs to restart. If empty, all nodes in the cluster will be restarted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response message from Wazuh after sending the restart request.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes for which restart was requested.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes where restart request failed.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the restart request was sent.",
			},
		},
	}
}

// Create (run) restart via PUT /cluster/restart
func resourceNodeRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	nodes := expandStringList(d.Get("nodes_list").([]interface{}))

	// Wazuh API expects nodes_list as query parameter, not in JSON body.
	// Body can be empty JSON to satisfy some HTTP clients.
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Build base URL
	urlStr := fmt.Sprintf("%s/cluster/restart", client.Endpoint)

	// Query parameters
	query := url.Values{}
	if len(nodes) > 0 {
		query.Set("nodes_list", joinComma(nodes))
	}
	if q := query.Encode(); q != "" {
		urlStr = urlStr + "?" + q
	}

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(bodyBytes))
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
		// Typicky 202 pro success, ale bereme jakýkoli 2xx jako OK
		return diag.Errorf("failed to restart nodes: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data struct {
			TotalAffected int `json:"total_affected_items"`
			TotalFailed   int `json:"total_failed_items"`
		} `json:"data"`
		Message string `json:"message"`
		Error   int    `json:"error"`
	}

	_ = json.Unmarshal(respBody, &result)

	if result.Error != 0 {
		return diag.Errorf("Wazuh API returned error while restarting nodes: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use timestamp-based ID so each restart is a unique apply
	d.SetId(fmt.Sprintf("node-restart-%d", time.Now().Unix()))

	return diags
}

// Read is a noop — this resource represents a one-time action.
func resourceNodeRestartRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// No-op for update/delete (cannot "undo" a restart request).
func resourceNodeRestartNoop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Import only re-attaches the existing ID to state (no API call).
func resourceNodeRestartImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
