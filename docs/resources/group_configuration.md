# 📦 **Resource Documentation: `wazuh_group_configuration`**

# wazuh_group_configuration

The `wazuh_group_configuration` resource allows you to manage and synchronize the **XML configuration (`agent.conf`)** for a specific agent group within a **Wazuh Manager** instance via its REST API.

This resource lets you **create**, **update**, **read**, and **import** group configurations directly from Terraform — making it easy to maintain Wazuh configuration as code.

---

## Example Usage

### Create or Update Wazuh Group Configuration

```hcl
resource "wazuh_group_configuration" "example" {
  group_id = "NewGroup_1"

  configuration_xml = <<EOF
<agent_config>
  <localfile>
    <log_format>syslog</log_format>
    <location>/var/log/auth.log</location>
  </localfile>
</agent_config>
EOF
}
```

### Import Existing Group Configuration

If a group already exists in Wazuh, you can import its configuration into Terraform state:

```bash
terraform import wazuh_group_configuration.example NewGroup_1
```

---

## Lifecycle & Behavior

* Creating or updating a `wazuh_group_configuration` resource triggers a `PUT /groups/{group_id}/configuration` request with the XML payload:

  ```xml
  <agent_config>
    <localfile>
      <log_format>syslog</log_format>
      <location>/var/log/auth.log</location>
    </localfile>
  </agent_config>
  ```

* Reading an existing configuration uses:

  ```
  GET /groups/{group_id}/configuration
  ```

* Deleting a configuration is **not supported** by the Wazuh API — Terraform will only remove the resource from its state.

Terraform automatically handles:

* Creating or updating (`PUT`),
* Refreshing state (`Read`),
* Importing existing configurations (`Import`).

> ⚠️ **Note:** Wazuh does not support partial updates for `agent.conf`. Each `PUT` request replaces the **entire configuration file**.

---

## Arguments Reference

| Name                | Type   | Required  | Description                                                                                                  |
| ------------------- | ------ | --------- | ------------------------------------------------------------------------------------------------------------ |
| `group_id`          | string | ✅ **Yes** | The ID (name) of the Wazuh group this configuration belongs to.                                              |
| `configuration_xml` | string | ✅ **Yes** | Full XML configuration content for the group (`agent.conf`). Must be valid XML as accepted by the Wazuh API. |

---

## Attributes Reference

| Name                | Description                                                       |
| ------------------- | ----------------------------------------------------------------- |
| `id`                | The unique ID of the configuration resource (same as `group_id`). |
| `configuration_xml` | The full XML content currently stored in Wazuh for the group.     |

## Notes

* The XML must include valid **Wazuh `agent.conf`** structure.
* For multi-line XML, use Terraform’s `<<EOF` heredoc syntax.
* The provider automatically sends the configuration as `application/xml` with JWT authorization.
