# 🚀 **Resource Documentation: `wazuh_agent`**

# wazuh_agent

The `wazuh_agent` resource allows you to **manage Wazuh agents** via the REST API.

It wraps:

* **Create** → `POST /agents/insert`
* **Read** → `GET /agents?agents_list=<id>`
* **Delete** → `DELETE /agents?agents_list=<id>&status=all&older_than=0s&purge=...`

and supports **Terraform import** for existing agents.

> 💡 This resource is focused on **agent registration / lifecycle on the manager** – it does not install binaries on the endpoint itself.

---

## Example Usage

### Minimal Agent Creation

```hcl
resource "wazuh_agent" "example" {
  name             = "NewHost_2"
  ip               = "10.0.10.11"
  purge_on_destroy = true
}
```

* `name` – Wazuh agent name
* `ip` – IP / IP/NET / `ANY` used during registration
* `purge_on_destroy = true` – permanently remove the agent from keystore on destroy

### Explicit Agent ID and Key

```hcl
resource "wazuh_agent" "fixed_id" {
  name             = "fixed-host"
  agent_id         = "123"
  ip               = "192.168.10.15"
  key              = "1abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghi64"
  purge_on_destroy = false
}
```

Here:

* You assign a specific `agent_id` and `key`,
* Wazuh will accept this identity if it’s consistent with the agent’s local `client.keys`.

### Import Existing Agent

If an agent already exists in Wazuh (e.g. registered via other tools), you can manage it with Terraform:

```bash
terraform import wazuh_agent.example 001
```

Then adjust `name`, `ip`, etc. in HCL to match reality (or let `terraform plan` show you the diff).

---

## Lifecycle & Behavior

### Create – `POST /agents/insert`

On `terraform apply`, the provider sends:

```http
POST /agents/insert
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "NewHost_2",
  "ip": "10.0.10.11",
  "id": "123",              // only if agent_id was set
  "key": "..."              // only if key was set
}
```

A typical response:

```json
{
  "data": {
    "id": "010",
    "key": "MDEwIHRlc3RfYWdlbnRfd2F..."
  },
  "error": 0
}
```

The resource then:

* sets `id` (Terraform ID) to the agent ID,
* sets `agent_id` from the response (if not already provided),
* sets `key` if returned,
* calls `GET /agents?agents_list=<id>` to populate computed attributes (`status`, `manager`, `version`, etc.).

---

### Read – `GET /agents?agents_list=<id>`

On refresh/plan, the provider runs:

```http
GET /agents?agents_list=<agent_id>
Authorization: Bearer <token>
```

If the agent exists, the response looks like:

```json
{
  "data": {
    "affected_items": [
      {
        "id": "001",
        "name": "NewHost_2",
        "ip": "10.0.10.11",
        "registerIP": "10.0.10.11",
        "status": "active",
        "manager": "wazuh.manager",
        "version": "Wazuh v4.14.0",
        "node_name": "node01"
      }
    ],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All selected agents information was returned",
  "error": 0
}
```

Values are mapped back to:

* `agent_id`, `name`, `ip`, `register_ip`, `status`, `manager`, `version`, `node_name`.

If the agent is not found (404 or empty `affected_items`), the provider clears the Terraform ID → resource is considered **destroyed externally**.

---

### Delete – `DELETE /agents`

On `terraform destroy`, the provider calls:

```http
DELETE /agents?agents_list=<agent_id>&status=all&older_than=0s&purge=<true|false>
Authorization: Bearer <token>
```

* `agents_list=<agent_id>` – selects this specific agent,
* `status=all` – regardless of current agent status,
* `older_than=0s` – don’t limit by last keep alive / registration time,
* `purge`:

  * `true` if `purge_on_destroy = true`,
  * `false` otherwise.

So you can control whether the agent is **soft-deleted** or **purged from keystore** via `purge_on_destroy`.

---

## Arguments Reference

| Name                              | Type   | Required    | ForceNew | Description                                                                                                                                                        |
|-----------------------------------|--------|-------------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`                            | string | ✅ **Yes**   | ✅ Yes   | Agent name in Wazuh. Changing this forces a new resource.                                                                                                          |
| `agent_id`                        | string | 🚫 optional | ✅ Yes   | Wazuh agent ID. If omitted, the API may assign one and it will be populated from the response. Changing this forces a new resource.                                |
| `ip`                              | string | 🚫 optional | ✅ Yes   | IP / IP/NET / `ANY`. If omitted, Wazuh will try to detect it. Changing this forces a new resource.                                                                 |
| `key`                             | string | 🚫 optional | ✅ Yes   | Shared key for communication with the manager. If omitted, Wazuh may generate/manage it separately. Changing this forces a new resource. Sensitive.                |
| `purge_on_destroy`               | bool   | 🚫 optional | ✅ Yes   | If `true`, the agent is permanently deleted from the key store on destroy (uses `purge=true`). Default `false`. Changing this forces a new resource.               |
| `force_enabled`                   | bool   | 🚫 optional | ✅ Yes   | Enable force insertion behavior. When `true`, the API may replace existing agents with the same name/ID/IP according to additional force conditions. Default `false`. |
| `force_disconnected_time_enabled` | bool   | 🚫 optional | ✅ Yes   | When using force, enable the `disconnected_time` condition. Default `true`. Changing this forces a new resource.                                                   |
| `force_disconnected_time_value`   | string | 🚫 optional | ✅ Yes   | Time the agent must have been disconnected to allow forced insertion (e.g. `\"30m\"`, `\"2h\"`, `\"7d\"`). Default `\"1h\"`. Changing this forces a new resource.  |
| `force_after_registration_time`   | string | 🚫 optional | ✅ Yes   | Time the agent must have been registered to allow forced insertion (e.g. `\"1h\"`, `\"2h\"`, `\"7d\"`). Default `\"1h\"`. Changing this forces a new resource.     |

> 🔁 Because the resource doesn’t support in-place `Update`, all configurable inputs are **ForceNew** – modifying any of them means Terraform will destroy + recreate the agent.

---

## Attributes Reference

| Name          | Type   | Description                                                                 |
| ------------- | ------ | --------------------------------------------------------------------------- |
| `id`          | string | Terraform resource ID – same as `agent_id`.                                 |
| `agent_id`    | string | Wazuh agent ID. Returned from API on create or import.                      |
| `name`        | string | Agent name in Wazuh.                                                        |
| `ip`          | string | IP configured for the agent.                                                |
| `register_ip` | string | IP used when registering the agent.                                         |
| `status`      | string | Agent status – e.g. `active`, `pending`, `never_connected`, `disconnected`. |
| `manager`     | string | Manager hostname where the agent is connected.                              |
| `version`     | string | Agent version string.                                                       |
| `node_name`   | string | Cluster node name the agent is attached to (if running in clustered setup). |
| `key`         | string | Shared key returned by Wazuh on create (if available). Sensitive.           |
