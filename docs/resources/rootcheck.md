# 🔍 **Resource Documentation: `wazuh_rootcheck`**

# wazuh_rootcheck

The `wazuh_rootcheck` resource manages **Wazuh rootcheck scans** for a specific agent via the `/rootcheck` API endpoints.

It lets you:

* **start a rootcheck scan** for an agent,
* **fetch a summary of rootcheck results** for that agent,
* **clear** the agent’s rootcheck database.

This resource is **per-agent** – one resource instance represents one Wazuh agent.

---

## Example Usage

```hcl
resource "wazuh_rootcheck" "agent_001" {
  agent_id = "001"
}

output "rootcheck_scan_message" {
  value = wazuh_rootcheck.agent_001.scan_message
}

output "rootcheck_results_message" {
  value = wazuh_rootcheck.agent_001.results_message
}

output "rootcheck_results_total_affected" {
  value = wazuh_rootcheck.agent_001.results_total_affected
}
```

---

## Lifecycle & Behavior

### Create – Run Rootcheck Scan

On `terraform apply`, the resource:

* Calls:

  ```http
  PUT /rootcheck?agents_list=001
  Content-Type: application/json
  ```

* Wazuh responds typically with:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 4,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "Rootcheck scan was restarted on returned agents",
    "error": 0
  }
  ```

* The resource stores:

  * `scan_message`
  * `scan_total_affected`
  * `scan_total_failed`
  * `last_scan_timestamp`

and uses `agent_id` as the Terraform resource `id`.

---

### Read – Get Rootcheck Results Summary

On refresh/plan, Terraform calls:

```http
GET /rootcheck/{agent_id}
```

and expects a response like:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 2,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All selected rootcheck information was returned",
  "error": 0
}
```

The resource updates:

* `results_message`
* `results_total_affected`
* `results_total_failed`

So you always have a **summary** of the agent’s rootcheck DB in Terraform state.

> ℹ️ Detailed per-item rootcheck data is not modeled as individual attributes – only the counts and message are tracked.

---

### Delete – Clear Rootcheck Database

On `terraform destroy`, the resource:

* Calls:

  ```http
  DELETE /rootcheck/{agent_id}
  ```

* Wazuh responds with:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 1,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "Rootcheck database was cleared on returned agents",
    "error": 0
  }
  ```

Terraform then removes the resource from state and the agent’s rootcheck DB is cleared.

> If the rootcheck data for that agent is already gone and Wazuh returns `404`, the provider treats it as success and still removes the resource from state.

---

## Arguments Reference

| Name       | Type   | Required  | Description                                                                                      |
| ---------- | ------ | --------- | ------------------------------------------------------------------------------------------------ |
| `agent_id` | string | ✅ **Yes** | Wazuh Agent ID (e.g. `"001"`). Must be at least 3 characters, matching the Wazuh AgentID format. |

---

## Attributes Reference

### Scan (Create) Result

| Name                  | Type   | Description                                                             |
| --------------------- | ------ | ----------------------------------------------------------------------- |
| `scan_message`        | string | Message returned by Wazuh when the rootcheck scan is started.           |
| `scan_total_affected` | int    | Number of agents for which rootcheck scan was successfully (re)started. |
| `scan_total_failed`   | int    | Number of agents where starting the rootcheck scan failed.              |
| `last_scan_timestamp` | string | UTC timestamp when the scan was triggered via Terraform.                |

### Results (Read) Summary

| Name                     | Type   | Description                                                         |
| ------------------------ | ------ | ------------------------------------------------------------------- |
| `results_message`        | string | Message returned when fetching rootcheck information for the agent. |
| `results_total_affected` | int    | Number of rootcheck items returned for this agent.                  |
| `results_total_failed`   | int    | Number of failed items when retrieving rootcheck information.       |

### General

| Name       | Type   | Description                        |
| ---------- | ------ | ---------------------------------- |
| `id`       | string | Resource ID (same as `agent_id`).  |
| `agent_id` | string | Agent ID managed by this resource. |
