# ⬆️ **Resource Documentation: `wazuh_agent_upgrade`**

# wazuh_agent_upgrade

The `wazuh_agent_upgrade` resource represents a **one-shot action** that creates **upgrade tasks for Wazuh agents** via the `/agents/upgrade` API endpoint.

It uses WPK packages from an online repository and lets you:

* upgrade **specific agents** by ID, or
* use the special keyword `all` to target **all agents**,
* optionally control repository, version, protocol, force behavior, and package type,
* get back the **task IDs** created per agent.

> This resource does **not** represent a persistent object.
> Each `apply` that creates this resource triggers a new upgrade request.

---

## Example Usage

### Upgrade Specific Agents to a Given Version

```hcl
resource "wazuh_agent_upgrade" "selected_agents" {
  agents_list     = ["001", "002"]
  upgrade_version = "4.14.0"
}

output "upgrade_message" {
  value = wazuh_agent_upgrade.selected_agents.message
}

output "upgrade_affected_items" {
  value = wazuh_agent_upgrade.selected_agents.affected_items
}
```

### Upgrade All Agents

```hcl
resource "wazuh_agent_upgrade" "all_agents" {
  agents_list = ["all"]
}

output "upgrade_all_message" {
  value = wazuh_agent_upgrade.all_agents.message
}
```

### Upgrade with Custom Options

```hcl
resource "wazuh_agent_upgrade" "custom_repo_upgrade" {
  agents_list     = ["001", "002", "003"]
  upgrade_version = "4.14.0"
  wpk_repo        = "packages.wazuh.com/4.x"
  package_type    = "deb"     # or "rpm"
  force           = true
  use_http        = false     # use HTTPS
}
```

---

## Lifecycle & Behavior

### Create – Create Upgrade Tasks

On `terraform apply`, the resource sends:

```http
PUT /agents/upgrade?agents_list=001,002&upgrade_version=4.14.0&...
Authorization: Bearer <token>
Content-Type: application/json
```

Depending on your configuration it might include query parameters such as:

* `agents_list=001,002` or `agents_list=all`
* `upgrade_version=4.14.0`
* `wpk_repo=packages.wazuh.com/4.x`
* `use_http=true|false`
* `force=true|false`
* `package_type=deb|rpm`

The body is a minimal JSON object (`{}`) for compatibility.

A typical successful response:

```json
{
  "data": {
    "affected_items": [
      {
        "agent": "001",
        "task_id": 1
      },
      {
        "agent": "002",
        "task_id": 2
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

* sets `message`,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the top-level `error` field,
* flattens `affected_items` into a list of `{ agent, task_id }`,
* sets `timestamp` to the current UTC time,
* sets `id` to a timestamp-based value (e.g. `20251108T211500Z`).

> ℹ️ A non-zero `error_code` doesn’t automatically fail the Terraform apply; HTTP status is used as the success check.
> `error_code` is exposed so you can inspect or react to partial failures in outputs / tooling.

---

### Read – No-op

On refresh/plan, `Read`:

* does **not** call the API again,
* simply keeps the values stored at `Create` time.

Task status tracking should be done via Wazuh’s task/status endpoints, not via this resource.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** cancel or roll back any upgrade tasks,
* only removes the action record from Terraform state.

The upgrade tasks themselves are handled by Wazuh.

---

## Arguments Reference

| Name              | Type         | Required    | Description                                                                                                                                                                                                                                  |
| ----------------- | ------------ | ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `agents_list`     | list(string) | ✅ **Yes**   | List of agent IDs to upgrade (e.g. `["001", "002"]`) **or** the keyword `"all"` to select all agents. Internally sent as `agents_list=001,002` or `agents_list=all`. Changes to this field force a new resource (and a new upgrade request). |
| `wpk_repo`        | string       | 🚫 optional | WPK repository URL/path. If omitted, Wazuh uses its default repository configuration.                                                                                                                                                        |
| `upgrade_version` | string       | 🚫 optional | Wazuh version to upgrade agents to (e.g. `"4.14.0"`). If omitted, Wazuh may infer a version depending on your configuration.                                                                                                                 |
| `use_http`        | bool         | 🚫 optional | Whether to use **HTTP** instead of **HTTPS** for repository access (`false` by default, i.e. HTTPS).                                                                                                                                         |
| `force`           | bool         | 🚫 optional | Force upgrade even if the agent seems to be already on the requested version.                                                                                                                                                                |
| `package_type`    | string       | 🚫 optional | Package type to use (`"rpm"` or `"deb"`). If omitted, the manager infers this automatically.                                                                                                                                                 |

> 🔁 All input arguments are **ForceNew** – changing any of them will cause Terraform to destroy/recreate the resource, which results in a new upgrade request.

---

## Attributes Reference

| Name             | Type         | Description                                                                                   |
| ---------------- | ------------ | --------------------------------------------------------------------------------------------- |
| `id`             | string       | Unique ID of this upgrade action, generated from the UTC timestamp (e.g. `20251108T211500Z`). |
| `agents_list`    | list(string) | Agents (IDs or `"all"`) targeted by this upgrade request.                                     |
| `message`        | string       | Human-readable message returned by the Wazuh API (e.g. `"All upgrade tasks were created"`).   |
| `total_affected` | int          | Number of agents for which upgrade tasks were created.                                        |
| `total_failed`   | int          | Number of agents where upgrade tasks could not be created.                                    |
| `error_code`     | int          | Raw `error` code from Wazuh API (`0` = success, `>0` = partial/failed).                       |
| `affected_items` | list(object) | List of objects `{ agent, task_id }` representing created upgrade tasks per agent.            |
| `timestamp`      | string       | UTC timestamp (RFC3339) when the upgrade request was sent via Terraform.                      |

`affected_items` has the structure:

```hcl
affected_items = [
  {
    agent   = "001"
    task_id = 1
  },
  {
    agent   = "002"
    task_id = 2
  },
]
```
