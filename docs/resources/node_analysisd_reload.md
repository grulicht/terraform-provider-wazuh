# ♻️ **Resource Documentation: `wazuh_node_analysisd_reload`**

# wazuh_node_analysisd_reload

The `wazuh_node_analysisd_reload` resource allows you to trigger a **reload of the `analysisd` process** on one or more Wazuh cluster nodes via the `/cluster/analysisd/reload` API endpoint.

This resource represents a **one-time action**: when applied, it sends a reload request to the target nodes.

> ⚠️ **Important:** This resource requires **Wazuh cluster mode** to be enabled.
> On non-cluster setups, the API may return an error indicating that the cluster is not running.

---

## Example Usage

### Reload `analysisd` on All Nodes

```hcl
resource "wazuh_node_analysisd_reload" "all_nodes" {
  # nodes_list omitted → all cluster nodes will be targeted
}
```

---

### Reload `analysisd` on Specific Nodes

```hcl
resource "wazuh_node_analysisd_reload" "some_nodes" {
  nodes_list = ["node01", "node02"]
}
```

---

### Optional Outputs

```hcl
output "analysisd_reload_message" {
  value = wazuh_node_analysisd_reload.all_nodes.message
}

output "analysisd_reload_total_affected" {
  value = wazuh_node_analysisd_reload.all_nodes.total_affected
}

output "analysisd_reload_total_failed" {
  value = wazuh_node_analysisd_reload.all_nodes.total_failed
}
```

---

## Lifecycle & Behavior

* Creating a `wazuh_node_analysisd_reload` resource triggers:

  ```http
  PUT /cluster/analysisd/reload[?nodes_list=node01,node02]
  Content-Type: application/json
  ```

  where:

  * `nodes_list` is passed as a **query parameter** (comma-separated list of node IDs),
  * if `nodes_list` is **not specified**, the reload request is sent to **all nodes** in the cluster.

* A typical successful response looks like:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 2,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "Reload request sent to all specified nodes",
    "error": 0
  }
  ```

* `terraform destroy` does **not revert** or roll back the reload.
  It only removes this one-time action from Terraform state.

> ⚠️ Every `terraform apply` with this resource present will send a **new reload request**.

---

## Arguments Reference

| Name         | Type         | Required    | Description                                                                                                              |
| ------------ | ------------ | ----------- | ------------------------------------------------------------------------------------------------------------------------ |
| `nodes_list` | list(string) | 🚫 optional | List of node IDs where `analysisd` should be reloaded. If omitted or empty, the reload request is sent to **all nodes**. |

---

## Attributes Reference

| Name             | Type   | Description                                                     |
| ---------------- | ------ | --------------------------------------------------------------- |
| `id`             | string | Internal Terraform resource ID (generated from timestamp).      |
| `message`        | string | Human-readable message from Wazuh about the reload request.     |
| `total_affected` | int    | Number of nodes where the reload request was sent successfully. |
| `total_failed`   | int    | Number of nodes where the reload request failed.                |
| `timestamp`      | string | UTC timestamp when the reload request was sent.                 |
