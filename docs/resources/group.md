# 📦 **Resource Documentation: `wazuh_group`**

# wazuh_group

The `wazuh_group` resource allows you to manage agent groups within a **Wazuh Manager** instance via its REST API.

Groups in Wazuh are used to logically organize and manage agents, configurations, and deployment rules.
This resource lets you **create**, **import**, **read**, and **delete** groups directly from Terraform.

---

## Example Usage

### Create Wazuh Group

```hcl
resource "wazuh_group" "example" {
  group_id = "NewGroup_1"
}
```

### Import Existing Group

If a group already exists in Wazuh, you can import it into Terraform state:

```bash
terraform import wazuh_group.example NewGroup_1
```

---

## Lifecycle & Behavior

* Creating a new `wazuh_group` resource triggers a `POST /groups` request with a payload like:

  ```json
  { "group_id": "NewGroup_1" }
  ```
* Reading an existing group uses:

  ```
  GET /groups?groups_list=NewGroup_1
  ```
* Deleting a group uses:

  ```
  DELETE /groups?groups_list=NewGroup_1
  ```

Terraform automatically handles:

* Refreshing state (`Read`),
* Importing existing groups (`Import`),
* Destroying groups (`Delete`).

> ⚠️ **Note:** Wazuh groups are immutable by ID. If you change the `group_id` value, Terraform will destroy and recreate the group.

---

## Arguments Reference

| Name       | Type   | Required  | Description                                                                                                                                         |
| ---------- | ------ | --------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `group_id` | string | ✅ **Yes** | The name of the Wazuh group. Must be 1–128 characters long and may include letters, numbers, `_`, `-`, and `.`. Names `.` and `..` are not allowed. |

---

## Attributes Reference

| Name | Description                                            |
| ---- | ------------------------------------------------------ |
| `id` | The unique ID of the Wazuh group (same as `group_id`). |

