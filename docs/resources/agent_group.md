# 🚀📦 **Resource Documentation: `wazuh_agent_group`**

# wazuh_agent_group

The `wazuh_agent_group` resource manages the relationship between **Wazuh agents** and **groups**.

It supports two modes:

1. **Single-agent mode**
   Manage membership of **one agent** in **one group** via:

   * `PUT /agents/{agent_id}/group/{group_id}`
   * `DELETE /agents/{agent_id}/group/{group_id}`

2. **Bulk mode**
   Manage membership of **multiple (or all) agents** in a group via:

   * `PUT /agents/group`
   * `DELETE /agents/group`

Use this resource when you want to:

* Assign an agent to a group (optionally as its **only** group),
* Assign a set of agents (or all agents) to a group,
* Remove an agent from a specific group,
* Remove a list of agents from a group,
* Capture the Wazuh API response (message, counts, error code, timestamp) as Terraform outputs.

> ⚠️ This resource represents an **action-like relationship**, not a full “desired state” of all agent-group assignments.
> It models *“this agent/these agents are (have been) assigned to group X via this operation”*.

---

## Example Usage

### 1. Single-Agent Mode – Assign Agent to Group

```hcl
resource "wazuh_agent_group" "agent_linux" {
  agent_id          = "001"
  group_id          = "linux"
  force_single_group = true
}

output "agent_linux_message" {
  value = wazuh_agent_group.agent_linux.message
}
```

This will:

* Call:

  ```http
  PUT /agents/001/group/linux?force_single_group=true
  ```

* Remove agent `001` from all existing groups,

* Assign it to `linux` as its only group.

---

### 2. Bulk Mode – Assign Selected Agents to Group

```hcl
resource "wazuh_agent_group" "bulk_app" {
  group_id    = "app-servers"
  agents_list = ["001", "002", "003"]
}
```

This will:

* Call:

  ```http
  PUT /agents/group?group_id=app-servers&agents_list=001,002,003
  ```

* Assign agents `001`, `002`, `003` to the `app-servers` group (without touching other groups, unless `force_single_group = true` is set).

---

### 3. Bulk Mode – Assign **All** Agents to a Group

If you omit `agent_id` **and** `agents_list`, Wazuh assigns **all agents** to the group:

```hcl
resource "wazuh_agent_group" "all_linux" {
  group_id          = "linux"
  force_single_group = false
}
```

This results in:

```http
PUT /agents/group?group_id=linux
```

Wazuh will:

* Assign all agents to `linux`,
* Keep existing group memberships (unless `force_single_group = true`).

---

### 4. Bulk Mode – Remove Agents from a Group

Bulk removal is done on **destroy** in bulk mode:

```hcl
resource "wazuh_agent_group" "temp_cleanup" {
  group_id    = "temp-group"
  agents_list = ["010", "011"]
}
```

* `terraform apply` assigns agents `010` and `011` to `temp-group`.
* `terraform destroy` calls:

  ```http
  DELETE /agents/group?group_id=temp-group&agents_list=010,011
  ```

---

## Lifecycle & Behavior

### Create

Depending on which arguments you set:

#### Single-Agent Mode *(agent_id set)*

```http
PUT /agents/{agent_id}/group/{group_id}?force_single_group=<bool>
```

* Used when `agent_id` is provided.
* Optional behavior:

  * `force_single_group = true` → remove agent from all other groups and assign only to `group_id`.
  * `force_single_group = false` → agent is added to `group_id` alongside existing groups.

Example response:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": [],
    "message": "All selected agents were assigned to group3",
    "error": 0
  }
}
```

#### Bulk Mode *(agent_id **not** set)*

```http
PUT /agents/group?group_id=<group_id>&agents_list=001,002&force_single_group=true
```

* If `agents_list` is set → only those agents are targeted.
* If `agents_list` is omitted → **all agents** are targeted.
* `force_single_group` behaves the same as in single-agent mode, but for multiple agents.

A typical bulk response:

```json
{
  "data": {
    "affected_items": ["001", "002"],
    "failed_items": [
      {
        "error": {
          "code": 1701,
          "message": "Agent does not exist",
          "remediation": "Please, use `GET /agents?select=id,name` to find all available agents"
        },
        "id": ["999"]
      }
    ],
    "total_affected_items": 2,
    "total_failed_items": 1,
    "message": "Some agents were not assigned to group2 and removed from the other groups",
    "error": 2
  }
}
```

Provider then:

* sets `message`,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the `error` field,
* sets `timestamp` to the current UTC time,
* sets `id` to:

  * `<agent_id>-<group_id>` in single-agent mode,
  * `<group_id>-<timestamp>` in bulk mode.

> ℹ️ Non-zero `error_code` does **not** automatically fail the Terraform apply – HTTP status code is used for success.
> You can inspect `error_code` and `total_failed` in outputs.

---

### Read

`Read` is a **no-op**:

* No additional Wazuh API request is made.
* The resource just keeps the last known values from `Create`/`Delete`.

To inspect current group membership in Wazuh, use the Wazuh API directly – this resource doesn’t try to reconcile the full set of groups for each agent.

---

### Delete

#### Single-Agent Mode

```http
DELETE /agents/{agent_id}/group/{group_id}
```

* Removes this specific agent from this specific group.
* Example response:

  ```json
  {
    "message": "Agent '004' removed from 'dmz'.",
    "error": 0
  }
  ```

#### Bulk Mode

```http
DELETE /agents/group?group_id=<group_id>&agents_list=001,002
```

* `agents_list` is **required** in bulk delete mode (when `agent_id` is not set).
* Wazuh removes the given agents from `group_id`.

In both cases, the provider:

* parses the response (single or bulk),
* updates `message`, `total_affected`, `total_failed`, `error_code`, `timestamp`,
* clears Terraform state (`id = ""`).

---

## Arguments Reference

| Name                 | Type         | Required    | ForceNew | Description                                                                                                                                                        |
| -------------------- | ------------ | ----------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `group_id`           | string       | ✅ **Yes**   | ✅ Yes    | Wazuh group ID (group name). Required in all modes.                                                                                                                |
| `agent_id`           | string       | 🚫 optional | ✅ Yes    | Wazuh agent ID (e.g. `"001"`). If set, single-agent endpoints `/agents/{agent_id}/group/{group_id}` are used. If omitted, bulk `/agents/group` endpoints are used. |
| `agents_list`        | list(string) | 🚫 optional | ✅ Yes    | List of agent IDs for bulk assign/remove via `/agents/group`. Ignored if `agent_id` is set. For bulk delete, this **must** be non-empty.                           |
| `force_single_group` | bool         | 🚫 optional | ✅ Yes    | If `true`, removes the agent(s) from all existing groups and assigns them only to `group_id`. Applies to both single-agent and bulk modes. Default: `false`.       |

> 🔁 All arguments are **ForceNew** – changing them will recreate the resource and thus trigger a new assign/remove operation.

---

## Attributes Reference

| Name             | Type         | Description                                                                                                             |
| ---------------- | ------------ | ----------------------------------------------------------------------------------------------------------------------- |
| `id`             | string       | Unique ID of this operation. In single-agent mode: `"<agent_id>-<group_id>"`. In bulk mode: `"<group_id>-<timestamp>"`. |
| `agent_id`       | string       | Agent ID used in single-agent mode (if any).                                                                            |
| `agents_list`    | list(string) | The list of agents targeted in bulk mode (if any).                                                                      |
| `group_id`       | string       | Wazuh group ID targeted by this operation.                                                                              |
| `message`        | string       | Human-readable message from Wazuh about the assignment/removal.                                                         |
| `total_affected` | int          | Number of items affected, as reported by the Wazuh API.                                                                 |
| `total_failed`   | int          | Number of failed items, as reported by the Wazuh API.                                                                   |
| `error_code`     | int          | Error value from Wazuh (`0` = success, `>0` = partial/failed).                                                          |
| `timestamp`      | string       | UTC timestamp (RFC3339) when the last operation was executed via Terraform.                                             |
