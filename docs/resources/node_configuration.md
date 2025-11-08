# ⚙️ **Resource Documentation: `wazuh_node_configuration`**# wazuh_node_configuration

The `wazuh_node_configuration` resource allows you to manage the **Wazuh node configuration (`ossec.conf`)** for a specific **cluster node** via the `/cluster/{node_id}/configuration` API endpoint.

> ⚠️ **Important:** This resource **requires Wazuh cluster mode to be enabled**.
> If the cluster is not running, the API returns error `3013`:
>
> > `Cluster is not running, it might be disabled in WAZUH_HOME/etc/ossec.conf`

Use this resource only in environments where Wazuh is properly configured as a cluster.

---

## Example Usage

### Manage Configuration of Node `node01`

```hcl
resource "wazuh_node_configuration" "node01" {
  node_id = "node01"

  configuration_xml = <<EOF
<ossec_config>
  <global>
    <email_notification>no</email_notification>
  </global>
</ossec_config>
EOF
}
```

This will **replace** the `ossec.conf` for the node `node01` in the cluster.

---

## Lifecycle & Behavior

* Creating or updating a `wazuh_node_configuration` resource triggers:

  ```http
  PUT /cluster/{node_id}/configuration
  Content-Type: application/octet-stream
  ```

  where the request body is the full `ossec.conf` XML.

* Reading the current configuration uses:

  ```http
  GET /cluster/{node_id}/configuration?raw=true
  ```

  and returns the raw XML content of the node’s configuration, which is stored in the `configuration_xml` attribute.

* Deleting this resource in Terraform **does not delete the configuration file** on the node.
  The Wazuh API does not provide a “delete config” endpoint. Terraform will only remove the resource from its state.

> ⚠️ **Cluster requirement:** If Wazuh cluster is disabled, `PUT /cluster/{node_id}/configuration` will fail with HTTP 400 and error `3013`.

---

## Arguments Reference

| Name                | Type   | Required  | Description                                                                                        |
| ------------------- | ------ | --------- | -------------------------------------------------------------------------------------------------- |
| `node_id`           | string | ✅ **Yes** | Name/ID of the cluster node whose configuration is managed (e.g. `node01`).                        |
| `configuration_xml` | string | ✅ **Yes** | Full XML content of the node’s `ossec.conf`. This replaces the existing configuration on the node. |

---

## Attributes Reference

| Name                | Description                                           |
| ------------------- | ----------------------------------------------------- |
| `id`                | Resource ID (same as `node_id`).                      |
| `node_id`           | Cluster node name managed by this resource.           |
| `configuration_xml` | Current XML configuration as read back from the node. |
