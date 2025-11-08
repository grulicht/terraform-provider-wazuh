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

// resourceNodeAnalysisdReload defines the Terraform resource schema and operations
// for reloading analysisd on Wazuh cluster nodes via PUT /cluster/analysisd/reload.
func resourceNodeAnalysisdReload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeAnalysisdReloadCreate,
		ReadContext:   resourceNodeAnalysisdReloadRead,
		UpdateContext: resourceNodeAnalysisdReloadNoop,
		DeleteContext: resourceNodeAnalysisdReloadNoop,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNodeAnalysisdReloadImport,
		},

		Schema: map[string]*schema.Schema{
			"nodes_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of node IDs to reload analysisd on. If empty, all nodes in the cluster will be targeted.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response message from Wazuh after sending the reload request.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes where the reload request was sent successfully.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes where the reload request failed.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the reload request was sent.",
			},
		},
	}
}

// Create (run) analysisd reload via PUT /cluster/analysisd/reload
func resourceNodeAnalysisdReloadCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	nodes := expandStringList(d.Get("nodes_list").([]interface{}))

	// Body can be empty JSON object; Wazuh uses query param for nodes_list.
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Base URL
	urlStr := fmt.Sprintf("%s/cluster/analysisd/reload", client.Endpoint)

	// Query params
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
		return diag.Errorf("failed to reload analysisd on nodes: status %d, body: %s", resp.StatusCode, string(respBody))
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
		return diag.Errorf("Wazuh API returned error while reloading analysisd: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Unique ID per execution
	d.SetId(fmt.Sprintf("analysisd-reload-%d", time.Now().Unix()))

	return diags
}

// Read is a noop — this is a one-time action resource.
func resourceNodeAnalysisdReloadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// No-op for update/delete
func resourceNodeAnalysisdReloadNoop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Import simply reattaches an existing ID to state.
func resourceNodeAnalysisdReloadImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
