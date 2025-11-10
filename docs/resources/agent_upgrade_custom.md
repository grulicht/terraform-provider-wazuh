# ⬆️ **Resource Documentation: `wazuh_agent_upgrade_custom`**

# wazuh_agent_upgrade_custom

The `wazuh_agent_upgrade_custom` resource represents a **one-shot action** that creates **custom upgrade tasks** for Wazuh agents using a **local WPK file**, via the `/agents/upgrade_custom` API endpoint.

Use it when you want to:

* upgrade agents using **your own WPK file** stored on the Wazuh manager,
* upgrade **specific agents** or **all agents** using the `all` keyword,
* get back the **task IDs** created per agent.

> This resource does **not** represent a persistent object.
> Each `apply` that creates this resource triggers exactly one custom upgrade request.

---

## Example Usage

### Upgrade Specific Agents from a Local WPK

```hcl
resource "wazuh_agent_upgrade_custom" "custom_wpk" {
  agents_list = ["001", "002"]
  file_path   = "/var/ossec/updates/custom-agent-4.14.0.wpk"
}

output "custom_upgrade_message" {
  value = wazuh_agent_upgrade_custom.custom_wpk.message
}

output "custom_upgrade_affected_items" {
  value = wazuh_agent_upgrade_custom.custom_wpk.affected_items
}
```

### Upgrade All Agents Using a Local WPK

```hcl
resource "wazuh_agent_upgrade_custom" "custom_wpk_all" {
  agents_list = ["all"]
  file_path   = "/var/ossec/updates/custom-agent-4.14.0.wpk"
}
```

### Specify a Custom Installer Script

```hcl
resource "wazuh_agent_upgrade_custom" "custom_installer" {
  agents_list = ["001", "002", "003"]
  file_path   = "/var/ossec/updates/custom-agent-4.14.0.wpk"
  installer   = "upgrade.sh" # or upgrade.bat for Windows agents
}
```

---

## Lifecycle & Behavior

### Create – Create Custom Upgrade Tasks

On `terraform apply`, the resource sends a request like:

```http
PUT /agents/upgrade_custom?agents_list=001,002&file_path=/var/ossec/updates/custom.wpk&installer=upgrade.sh
Authorization: Bearer <token>
Content-Type: application/json
```

depending on the arguments you provide:

* `agents_list=001,002` or `agents_list=all`
* `file_path=/var/ossec/...`
* `installer=upgrade.sh` (optional)

The request body is a minimal JSON object (`{}`) for compatibility.

A typical successful response:

```json
{
  "data": {
    "affected_items": [
      {
        "agent": "005",
        "task_id": 3
      },
      {
        "agent": "006",
        "task_id": 4
      }
    ],
    "total_affected_items": 2,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All upgrade tasks were created",
  "error": 0
}
```

The resource then:

* sets `message` from the response,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the `error` field,
* flattens `affected_items` to a list of objects `{ agent, task_id }`,
* sets `timestamp` to the current UTC time,
* sets `id` to a timestamp-based value (e.g. `20251108T212300Z`).

> ℹ️ Non-zero `error_code` does **not** automatically fail the Terraform apply; success/failure is based on HTTP status.
> `error_code` is exposed for inspection and troubleshooting.

---

### Read – No-op

On refresh/plan, `Read`:

* does **not** query Wazuh again,
* keeps the values obtained at `Create` time.

It does not track ongoing task status – use Wazuh task APIs or dashboard for that.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** cancel upgrade tasks,
* only removes this action record from Terraform state.

All upgrade tasks remain managed by Wazuh.

---

## Arguments Reference

| Name          | Type         | Required    | Description                                                                                                                                               |
| ------------- | ------------ | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `agents_list` | list(string) | ✅ **Yes**   | List of agent IDs to upgrade (e.g. `["001", "002"]`) **or** the keyword `"all"` to target all agents. Sent as `agents_list=001,002` or `agents_list=all`. |
| `file_path`   | string       | ✅ **Yes**   | Full path to the WPK file on the Wazuh manager (must be inside the Wazuh installation directory, typically `/var/ossec`).                                 |
| `installer`   | string       | 🚫 optional | Installation script to use (e.g. `upgrade.sh` or `upgrade.bat`). If omitted, Wazuh uses its default script for the platform.                              |

> 🔁 All input fields are **ForceNew** – changing any of them will cause Terraform to destroy/recreate the resource and send a new custom upgrade request.

---

## Attributes Reference

| Name             | Type         | Description                                                                                       |
| ---------------- | ------------ | ------------------------------------------------------------------------------------------------- |
| `id`             | string       | Unique ID of this custom upgrade request, generated from UTC timestamp (e.g. `20251108T212300Z`). |
| `agents_list`    | list(string) | Agents targeted by this upgrade request (IDs or `"all"`).                                         |
| `file_path`      | string       | WPK file path used in the request.                                                                |
| `installer`      | string       | Installer script used (if set).                                                                   |
| `message`        | string       | Human-readable message returned by Wazuh (e.g. `"All upgrade tasks were created"`).               |
| `total_affected` | int          | Number of agents for which upgrade tasks were created.                                            |
| `total_failed`   | int          | Number of agents where tasks could not be created.                                                |
| `error_code`     | int          | Raw `error` code from the Wazuh API (`0` = success, `>0` = partial/failed).                       |
| `affected_items` | list(object) | List of `{ agent, task_id }` pairs representing created upgrade tasks per agent.                  |
| `timestamp`      | string       | UTC timestamp (RFC3339) when the request was sent.                                                |

`affected_items` structure:

```hcl
affected_items = [
  {
    agent   = "005"
    task_id = 3
  },
  {
    agent   = "006"
    task_id = 4
  },
]
```
