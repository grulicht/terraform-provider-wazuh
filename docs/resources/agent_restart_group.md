# 🔄 **Resource Documentation: `wazuh_agent_restart_group`**

# wazuh_agent_restart_group

The `wazuh_agent_restart_group` resource represents a **one-shot action** that restarts **all agents belonging to a specific Wazuh group** using the:

```http
PUT /agents/group/{group_id}/restart
```

API endpoint.

Use it when you want to:

* Restart all agents in a **single logical group** (e.g. `linux`, `windows`, `production`),
* Drive that restart from Terraform,
* Capture the Wazuh API response (message, counts, error code, timestamp) as outputs.

> This is an **action resource** – not a long-lived configuration object.
> Each `apply` triggers a new restart command for the target group.

---

## Example Usage

### Restart All Agents in a Group

```hcl
resource "wazuh_agent_restart_group" "linux_agents" {
  group_id = "linux"
}

output "restart_linux_message" {
  value = wazuh_agent_restart_group.linux_agents.message
}

output "restart_linux_total_affected" {
  value = wazuh_agent_restart_group.linux_agents.total_affected
}
```

---

## Lifecycle & Behavior

### Create – Send Restart Command

On `terraform apply`, this resource:

1. Builds the URL:

   ```http
   PUT /agents/group/{group_id}/restart
   ```

   For example, with `group_id = "linux"`:

   ```http
   PUT /agents/group/linux/restart
   ```

2. Sends a minimal JSON body (e.g. `{}`) along with the `Authorization: Bearer <token>` header.

A typical response from Wazuh looks like:

```json
{
  "data": {
    "affected_items": [
      "002",
      "003",
      "005"
    ],
    "total_affected_items": 3,
    "total_failed_items": 0,
    "failed_items": [],
    "message": "Restart command was sent to all agents",
    "error": 0
  }
}
```

The resource then:

* sets `message` from `data.message`,
* sets `total_affected` from `data.total_affected_items`,
* sets `total_failed` from `data.total_failed_items`,
* sets `error_code` from `data.error`,
* sets `timestamp` to the current UTC time (RFC3339),
* sets `id` to a unique value such as `<group_id>-20251108T205500Z`.

> ℹ️ The provider uses **HTTP status code** to decide success/failure.
> The Wazuh `error` field is exposed via `error_code` so you can inspect success/partial failure in outputs or further logic.

---

### Read – No-op

On refresh/plan (`Read`):

* No additional Wazuh API call is made.
* The resource simply keeps the values that were stored at create time.

This is intentional: the restart is a **point-in-time action**, not a stateful object to reconcile.

---

### Delete – State Only

On `terraform destroy`, this resource:

* does **not** call any Wazuh API,
* only removes this “restart action” from the Terraform state.

No agents are restarted again, and nothing is rolled back – the action has already happened during `Create`.

---

## Arguments Reference

| Name       | Type   | Required  | ForceNew | Description                                                                                             |
| ---------- | ------ | --------- | -------- | ------------------------------------------------------------------------------------------------------- |
| `group_id` | string | ✅ **Yes** | ✅ Yes    | Wazuh group ID (group name) whose agents will be restarted. Changing this forces a new action/resource. |

> 🔁 `group_id` is **ForceNew** – changing it will destroy/recreate the resource, triggering a new restart command for the new group.

---

## Attributes Reference

| Name             | Type   | Description                                                                                           |
| ---------------- | ------ | ----------------------------------------------------------------------------------------------------- |
| `id`             | string | Unique ID of this restart action, typically `<group_id>-<timestamp>` (e.g. `linux-20251108T205500Z`). |
| `group_id`       | string | The Wazuh group whose agents were targeted by this restart action.                                    |
| `message`        | string | Human-readable message returned by Wazuh after sending the restart command.                           |
| `total_affected` | int    | Number of agents for which the restart command was processed.                                         |
| `total_failed`   | int    | Number of agents where the restart command failed.                                                    |
| `error_code`     | int    | Raw `error` value from the Wazuh API response (`0` = success, `>0` = partial/failed states).          |
| `timestamp`      | string | UTC timestamp (RFC3339) when the restart request was sent via Terraform.                              |
