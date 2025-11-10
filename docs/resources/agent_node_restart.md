# 🔄 **Resource Documentation: `wazuh_agent_node_restart`**

# wazuh_agent_node_restart

The `wazuh_agent_node_restart` resource represents a **one-shot action** that restarts **all agents belonging to a specific Wazuh cluster node**, using the `/agents/node/{node_id}/restart` API endpoint.

> ⚠️ **Important:** This endpoint is only meaningful in a **multi-node Wazuh cluster**.
> In a single-node (non-clustered) deployment, this API may not be available or may not behave as expected.

Use this when you want to:

* restart all agents **attached to a particular node**,
* orchestrate node-local agent restarts after config / routing changes,
* capture the response (message, counts, error code, timestamp) as Terraform outputs.

This is an **action resource**, not a persistent object: each `apply` that creates it triggers exactly one restart command.

---

## Example Usage

```hcl
resource "wazuh_agent_node_restart" "node01_agents" {
  node_id = "node01"
}

output "node01_agents_restart_message" {
  value = wazuh_agent_node_restart.node01_agents.message
}

output "node01_agents_restart_total_affected" {
  value = wazuh_agent_node_restart.node01_agents.total_affected
}
```

---

## Lifecycle & Behavior

### Create – Restart Agents on a Node

On `terraform apply`, the resource:

* Builds and sends:

  ```http
  PUT /agents/node/{node_id}/restart
  Authorization: Bearer <token>
  Content-Type: application/json
  ```

  where `{node_id}` is the `node_id` attribute, e.g. `node01`.

* Sends an empty/minimal JSON body (e.g. `{}`) for compatibility.

A typical response:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 3,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "Restart command was sent to all agents",
  "error": 0
}
```

The resource then:

* sets `message` from the API response,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the `error` field,
* sets `timestamp` to the current UTC time,
* sets `id` to a combination of `node_id` and a timestamp, e.g. `node01-20251108T210900Z`.

> ℹ️ Non-zero `error_code` does **not automatically fail** the Terraform apply; the HTTP status code controls success/failure.
> The `error_code` is exposed for you to inspect in outputs/logs.

---

### Read – No-op

On refresh/plan, `Read`:

* does **not** query the Wazuh API again,
* keeps the values that were stored during `Create`.

This is intentional – a node-level agent restart is a one-time action, not a configuration object to be synchronized.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** call any API,
* only removes this restart action from the Terraform state.

The agents are **not** restarted again nor “un-restarted” – the action already happened when the resource was created.

---

## Arguments Reference

| Name      | Type   | Required  | Description                                                                                                                                                                                                 |
| --------- | ------ | --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `node_id` | string | ✅ **Yes** | Cluster node name whose agents should be restarted (e.g. `node01`). Must correspond to an existing node in a **Wazuh cluster**. Changing this value forces a new resource (and thus a new restart command). |

> 🔁 `node_id` is **ForceNew** – any change causes Terraform to destroy/recreate the resource, which triggers a new restart action for the new node.

---

## Attributes Reference

| Name             | Type   | Description                                                                                                       |
| ---------------- | ------ | ----------------------------------------------------------------------------------------------------------------- |
| `id`             | string | Unique ID of this restart action, typically `node_id` combined with a timestamp (e.g. `node01-20251108T210900Z`). |
| `node_id`        | string | Name of the cluster node whose agents were restarted.                                                             |
| `message`        | string | Human-readable message returned by Wazuh after sending the restart command.                                       |
| `total_affected` | int    | Number of agents on that node for which the restart command was processed.                                        |
| `total_failed`   | int    | Number of agents on that node where the restart command failed.                                                   |
| `error_code`     | int    | Raw `error` code from the Wazuh API (`0` = success, `>0` = partial/failed).                                       |
| `timestamp`      | string | UTC timestamp (RFC3339) when the restart request was sent via Terraform.                                          |
