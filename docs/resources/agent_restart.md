# 🔄 **Resource Documentation: `wazuh_agent_restart`**

# wazuh_agent_restart

The `wazuh_agent_restart` resource represents a **one-shot action** that restarts:

* **all Wazuh agents**,
* a specific **list of agents**, or
* a **single agent** via its dedicated endpoint

using the Wazuh `/agents/restart` and `/agents/{agent_id}/restart` API endpoints.

Use it when you want to:

* orchestrate agent restarts from Terraform,
* optionally target only selected agents,
* capture the API response (message, counts, error code, timestamp) as outputs.

> This is an **action resource** – not a long-lived thing.
> Each `apply` with this resource triggers a **new restart command**.

---

## Example Usage

### Restart All Agents

```hcl
resource "wazuh_agent_restart" "all_agents" {}

output "agent_restart_all_message" {
  value = wazuh_agent_restart.all_agents.message
}
```

### Restart a Single Agent

```hcl
resource "wazuh_agent_restart" "single_agent" {
  agents_list = ["001"]
}

output "agent_restart_single_message" {
  value = wazuh_agent_restart.single_agent.message
}
```

> For a single agent, the provider uses the **dedicated endpoint**:
>
> ```http
> PUT /agents/001/restart
> ```

### Restart Specific Agents (Multiple)

```hcl
resource "wazuh_agent_restart" "some_agents" {
  agents_list = ["001", "002"]
}

output "agent_restart_some_message" {
  value = wazuh_agent_restart.some_agents.message
}

output "agent_restart_some_total_affected" {
  value = wazuh_agent_restart.some_agents.total_affected
}

output "agent_restart_some_error_code" {
  value = wazuh_agent_restart.some_agents.error_code
}
```

---

## Lifecycle & Behavior

### Create – Send Restart Command

On `terraform apply`, this resource decides which Wazuh endpoint to call based on `agents_list`:

* If `agents_list` is **not set or empty** → restart **all agents**:

  ```http
  PUT /agents/restart
  ```

* If `agents_list` contains **exactly one ID**, e.g. `["001"]` → restart a **single agent**:

  ```http
  PUT /agents/001/restart
  ```

* If `agents_list` contains **multiple IDs**, e.g. `["001", "002"]` → restart **only those agents**:

  ```http
  PUT /agents/restart?agents_list=001,002
  ```

A minimal JSON body (e.g. `{}`) is sent with the request when required by the Wazuh API.

A typical response looks like:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 8,
    "total_failed_items": 3,
    "failed_items": []
  },
  "message": "Restart command was not sent to some agents",
  "error": 2
}
```

The resource then:

* sets `message` from the API response,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the top-level `error` field,
* sets `timestamp` to current UTC time,
* sets `id` to a unique timestamp-based string (e.g. `20251108T205500Z`).

> ℹ️ **Non-zero `error_code` does *not* automatically fail the Terraform apply** – the provider treats the HTTP status code as the success criterion and exposes `error_code` so you can inspect it in outputs.

---

### Read – No-op

On refresh/plan (`Read`):

* There is **no additional API call**.
* The resource simply keeps the stored values from the last `Create`.

This is by design – the restart command is a one-time action, not a persistent state to reconcile.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** call any Wazuh API again,
* only removes this “restart action” from Terraform state.

Agents are **not** restarted again and nothing is rolled back – the action already happened at create time.

---

## Arguments Reference

| Name          | Type         | Required    | Description                                                                                                                                                                                                                                                                         |
| ------------- | ------------ | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `agents_list` | list(string) | 🚫 optional | Optional list of Wazuh agent IDs to restart (e.g. `["001", "002"]`). If omitted or empty, the provider restarts **all agents**. If one ID is provided, the provider uses `/agents/{agent_id}/restart`. Changes to this field force a new resource (and thus a new restart command). |

> 🔁 `agents_list` is marked as **ForceNew** – changing it will cause Terraform to destroy/recreate the resource, which triggers a new restart command with the new target set.

---

## Attributes Reference

| Name             | Type         | Description                                                                                         |
| ---------------- | ------------ | --------------------------------------------------------------------------------------------------- |
| `id`             | string       | Unique ID of this restart action, generated from the UTC timestamp (e.g. `20251108T205500Z`).       |
| `agents_list`    | list(string) | The list of agents that this restart action targeted (if any).                                      |
| `message`        | string       | Human-readable message returned by Wazuh after sending the restart command.                         |
| `total_affected` | int          | Number of agents for which the restart command was processed.                                       |
| `total_failed`   | int          | Number of agents where the restart command failed.                                                  |
| `error_code`     | int          | Raw `error` code from the Wazuh API response (`0` = success, `>0` indicates partial/failed states). |
| `timestamp`      | string       | UTC timestamp (RFC3339) when the restart request was sent via Terraform.                            |
