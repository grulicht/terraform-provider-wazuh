# 📊 **Resource Documentation: `wazuh_syscheck`**

# wazuh_syscheck

The `wazuh_syscheck` resource manages **Wazuh File Integrity Monitoring (FIM / Syscheck)** for a specific agent via the `/syscheck` API endpoints.

With this resource you can:

* **start a Syscheck/FIM scan** for a given agent,
* **retrieve a summary of FIM findings** for that agent,
* **clear Syscheck database** for older agents (where supported).

This resource is **per-agent** – one instance represents Syscheck state for one Wazuh agent.

---

## Example Usage

```hcl
resource "wazuh_syscheck" "agent_001" {
  agent_id = "001"
}

output "syscheck_scan_message" {
  value = wazuh_syscheck.agent_001.scan_message
}

output "syscheck_results_message" {
  value = wazuh_syscheck.agent_001.results_message
}

output "syscheck_results_total_affected" {
  value = wazuh_syscheck.agent_001.results_total_affected
}
```

---

## Lifecycle & Behavior

### Create – Run FIM Scan

On `terraform apply`, the resource:

* Calls:

  ```http
  PUT /syscheck?agents_list=001
  Content-Type: application/json
  ```

* A typical response:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 4,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "Syscheck scan was restarted on returned agents",
    "error": 0
  }
  ```

* The resource then stores:

  * `scan_message`
  * `scan_total_affected`
  * `scan_total_failed`
  * `last_scan_timestamp`

and uses `agent_id` as the Terraform resource `id`.

---

### Read – Get FIM Findings Summary

On refresh/plan (`terraform plan`, `terraform refresh`), Terraform calls:

```http
GET /syscheck/{agent_id}
```

with responses like:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 1433,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "FIM findings of the agent were returned",
  "error": 0
}
```

The provider then updates:

* `results_message`
* `results_total_affected`
* `results_total_failed`

> ℹ️ Only summary counts and message are stored – detailed per-file findings are not broken out into Terraform attributes.

If the API returns `404` (no data for this agent), the resource is removed from state.

---

### Delete – Clear FIM Results (Where Supported)

On `terraform destroy`, the resource executes:

```http
DELETE /syscheck/{agent_id}
```

which clears the Syscheck database for that agent (for agents **< 3.12.0**, according to Wazuh docs).
For newer agents, this may be a no-op on the server side, but should still return success.

A typical response:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "Syscheck database was cleared on returned agents",
  "error": 0
}
```

> If Wazuh returns `404` (already cleared / no DB), the provider treats this as success and removes the resource from Terraform state.

---

## Arguments Reference

| Name       | Type   | Required  | Description                                                                            |
| ---------- | ------ | --------- | -------------------------------------------------------------------------------------- |
| `agent_id` | string | ✅ **Yes** | Wazuh Agent ID (e.g. `"001"`). Must follow standard Wazuh agent ID format (3+ digits). |

---

## Attributes Reference

### Scan (Create) Result

| Name                  | Type   | Description                                                       |
| --------------------- | ------ | ----------------------------------------------------------------- |
| `scan_message`        | string | Message returned by Wazuh when the syscheck scan is started.      |
| `scan_total_affected` | int    | Number of agents for which the syscheck scan was restarted.       |
| `scan_total_failed`   | int    | Number of agents where starting the syscheck scan failed.         |
| `last_scan_timestamp` | string | UTC timestamp when the syscheck scan was triggered via Terraform. |

### Results (Read) Summary

| Name                     | Type   | Description                                                 |
| ------------------------ | ------ | ----------------------------------------------------------- |
| `results_message`        | string | Message returned when fetching FIM findings for the agent.  |
| `results_total_affected` | int    | Total number of FIM findings/items returned for this agent. |
| `results_total_failed`   | int    | Number of failed items when retrieving FIM findings.        |

### General

| Name       | Type   | Description                        |
| ---------- | ------ | ---------------------------------- |
| `id`       | string | Resource ID (same as `agent_id`).  |
| `agent_id` | string | Agent ID managed by this resource. |
