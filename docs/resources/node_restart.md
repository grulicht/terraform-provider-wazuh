# 🔄 **Resource Documentation: `wazuh_node_restart`**

# wazuh_node_restart

The `wazuh_node_restart` resource allows you to trigger a **restart of one or more Wazuh cluster nodes** via the `/cluster/restart` API endpoint.

This resource represents a **one-time action**: when you apply it, Terraform sends a restart request to the specified nodes (or all nodes if none are specified).

> ⚠️ **Important:** This resource requires **Wazuh cluster mode** to be enabled.
> On non-cluster setups, the API may return an error indicating that the cluster is not running.

---

## Example Usage

### Restart All Nodes in the Cluster

```hcl
resource "wazuh_node_restart" "all_nodes" {
  # nodes_list is omitted → restart all nodes
}
```

---

### Restart Specific Nodes

```hcl
resource "wazuh_node_restart" "some_nodes" {
  nodes_list = ["node01", "node02"]
}
```

---

### Optional Outputs

```hcl
output "node_restart_message" {
  value = wazuh_node_restart.all_nodes.message
}

output "node_restart_total_affected" {
  value = wazuh_node_restart.all_nodes.total_affected
}

output "node_restart_total_failed" {
  value = wazuh_node_restart.all_nodes.total_failed
}
```

---

## Lifecycle & Behavior

* Creating a `wazuh_node_restart` resource triggers:

  ```http
  PUT /cluster/restart[?nodes_list=node01,node02]
  Content-Type: application/json
  ```

  where:

  * `nodes_list` is provided as a **query parameter** (comma-separated),
  * if `nodes_list` is not provided, the restart request is sent to **all nodes**.

* The response on success is typically:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 3,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "Restart request sent to all specified nodes",
    "error": 0
  }
  ```

* `terraform destroy` does **not** roll back or “undo” the restart.
  It only removes the record of this action from the Terraform state.

> ⚠️ Each `apply` with this resource will send a **new restart request**. Treat it like a “run command” action, not a persistent object.

---

## Arguments Reference

| Name         | Type         | Required    | Description                                                                                                    |
| ------------ | ------------ | ----------- | -------------------------------------------------------------------------------------------------------------- |
| `nodes_list` | list(string) | 🚫 optional | List of node IDs to restart. If omitted or empty, the restart request is sent to **all nodes** in the cluster. |

---

## Attributes Reference

| Name             | Type   | Description                                                          |
| ---------------- | ------ | -------------------------------------------------------------------- |
| `id`             | string | Internal Terraform resource ID (generated from timestamp).           |
| `message`        | string | Human-readable message from Wazuh about the restart request.         |
| `total_affected` | int    | Number of nodes for which the restart request was successfully sent. |
| `total_failed`   | int    | Number of nodes where the restart request failed.                    |
| `timestamp`      | string | UTC timestamp when the restart request was initiated.                |
