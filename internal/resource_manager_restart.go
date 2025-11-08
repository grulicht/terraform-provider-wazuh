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

// resourceManagerRestart models an action that restarts the Wazuh manager
// via PUT /manager/restart.
func resourceManagerRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceManagerRestartCreate,
		ReadContext:   resourceManagerRestartRead,
		// No Update, it's a one-shot action
		DeleteContext: resourceManagerRestartDelete,

		// No importer – importing a past restart action doesn't make sense
		Schema: map[string]*schema.Schema{
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Message returned by Wazuh after sending the restart request.",
			},
			"total_affected": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of managers/nodes affected by the restart request.",
			},
			"total_failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of managers/nodes where the restart request failed.",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UTC timestamp when the restart request was sent.",
			},
		},
	}
}

// Create: send restart request via PUT /manager/restart
func resourceManagerRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	var diags diag.Diagnostics

	urlStr := fmt.Sprintf("%s/manager/restart", client.Endpoint)

	// some Wazuh endpoints accept empty JSON; we send {} to be safe
	bodyBytes, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", urlStr, bytes.NewBuffer(bodyBytes))
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
		return diag.Errorf("failed to restart manager: status %d, body: %s", resp.StatusCode, string(respBody))
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
		return diag.Errorf("Wazuh API returned error while restarting manager: %s", string(respBody))
	}

	_ = d.Set("message", result.Message)
	_ = d.Set("total_affected", result.Data.TotalAffected)
	_ = d.Set("total_failed", result.Data.TotalFailed)
	_ = d.Set("timestamp", time.Now().UTC().Format(time.RFC3339))

	// Use timestamp as ID to make each restart action unique
	d.SetId(time.Now().UTC().Format("20060102T150405Z"))

	return diags
}

// Read: nothing to refresh; restart is a one-shot action
func resourceManagerRestartRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// no-op, we just keep the stored values
	return diags
}

// Delete: no API call – just forget this restart action from state
func resourceManagerRestartDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
