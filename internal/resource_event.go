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

// resourceEvent defines the Terraform resource schema and CRUD operations
// for ingesting security events into Wazuh (POST /events).
func resourceEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventCreate,
		ReadContext:   resourceEventRead,
		UpdateContext: resourceEventNoop,
		DeleteContext: resourceEventNoop,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEventImport,
		},

		Schema: map[string]*schema.Schema{
			"events": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of events to ingest. Each element is a string (plain text or JSON string). Max 100 events per request.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description returned by Wazuh after ingesting events.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of events successfully forwarded to analysisd.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of events that failed to be forwarded.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the events were ingested.",
			},
		},
	}
}

// Create (ingest) events via POST /events
func resourceEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	events := expandStringList(d.Get("events").([]interface{}))

	payload := map[string]interface{}{
		"events": events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("%s/events", client.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
		return diag.Errorf("failed to ingest events: status %d, body: %s", resp.StatusCode, string(respBody))
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
		return diag.Errorf("Wazuh API returned error while ingesting events: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use timestamp as unique ID so each ingestion is a unique resource instance
	d.SetId(fmt.Sprintf("events-%d", time.Now().Unix()))

	return diags
}

// Read is a noop: there is nothing to fetch back from the API for a past ingestion.
func resourceEventRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Noop for update/delete: ingested events cannot be undone.
func resourceEventNoop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

// Import simply attaches an existing ID to state; no API call.
func resourceEventImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
