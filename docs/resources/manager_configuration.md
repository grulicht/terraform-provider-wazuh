# ⚙️ **Resource Documentation: `wazuh_manager_configuration`**

# wazuh_manager_configuration

The `wazuh_manager_configuration` resource manages the **global Wazuh manager configuration** (`ossec.conf`) via the `/manager/configuration` API endpoint.

It allows you to:

* **replace the entire Wazuh manager configuration** (ossec.conf) with Terraform-managed XML,
* **read back the current manager configuration** as plain XML.

> ⚠️ There is **no “delete configuration”** in the Wazuh API.
> Destroying this resource only removes it from **Terraform state** – it does **not** reset `ossec.conf` on the Wazuh side.

---

## Example Usage

```hcl
resource "wazuh_manager_configuration" "manager" {
  configuration_xml = <<-EOT
    <ossec_config>
      ...
    </ossec_config>
  EOT
}

output "manager_config_message" {
  value = wazuh_manager_configuration.manager.message
}
```

* `terraform apply` uploads the provided `configuration_xml` as the new manager configuration.
* `terraform plan` / `terraform refresh` fetches the current configuration from Wazuh and syncs it back into state.

---

## Lifecycle & Behavior

### Create / Update – Replace Manager Configuration

On `terraform apply` (Create/Update), the resource calls:

```http
PUT /manager/configuration
Content-Type: application/octet-stream
```

with the raw XML from `configuration_xml` as the request body.

A successful response looks like:

```json
{
  "message": "Configuration was successfully updated",
  "error": 0
}
```

The resource:

* stores the `message` from the API,
* sets `last_updated_timestamp` to the current UTC timestamp,
* sets a fixed `id` of `"manager"`.

> ⚠️ This operation **replaces** the entire manager configuration file.
> Make sure your `configuration_xml` is a complete, valid `ossec.conf`, not a fragment.

---

### Read – Get Current Manager Configuration

On refresh/plan, Terraform calls:

```http
GET /manager/configuration?raw=true
```

* With `raw=true`, Wazuh returns the current `ossec.conf` content as **plain XML**.
* The provider updates:

  * `configuration_xml` – with the exact content returned by the API.

If Wazuh returns `404` (which would be unusual for a manager config), the resource is removed from Terraform state.

---

### Delete – State Only

On `terraform destroy`, the resource:

* **does not call any API** (Wazuh has no “reset configuration” endpoint),
* simply removes the resource from Terraform state.

The manager configuration on the Wazuh side remains unchanged.

> If you want to “revert” a configuration, do it by **applying a different `configuration_xml`**, not by destroying this resource.

---

## Arguments Reference

| Name                | Type   | Required  | Description                                                                |
| ------------------- | ------ | --------- | -------------------------------------------------------------------------- |
| `configuration_xml` | string | ✅ **Yes** | Full `ossec.conf` XML content to apply as the Wazuh manager configuration. |

---

## Attributes Reference

| Name                     | Type   | Description                                                                      |
| ------------------------ | ------ | -------------------------------------------------------------------------------- |
| `id`                     | string | Resource ID. Always set to `"manager"` for this resource.                        |
| `configuration_xml`      | string | Current manager configuration as last read from Wazuh (or applied by Terraform). |
| `message`                | string | Message returned by Wazuh after a successful `PUT /manager/configuration` call.  |
| `last_updated_timestamp` | string | UTC timestamp when the manager configuration was last updated via this resource. |

---

## Importing Existing Configuration

If you already have a manager configuration in place and want to bring it under Terraform control:

```bash
terraform import wazuh_manager_configuration.manager manager
```

> The import ID is arbitrary but normalized to `"manager"` internally.
> After import, run `terraform refresh` to populate `configuration_xml` from Wazuh.

You can then:

* format / restructure the XML in your `.tf` files,
* or replace it entirely with your desired configuration and re-apply.
