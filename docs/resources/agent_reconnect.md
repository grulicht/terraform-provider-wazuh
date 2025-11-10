# 🔗 **Resource Documentation: `wazuh_agent_reconnect`**

# wazuh_agent_reconnect

The `wazuh_agent_reconnect` resource is a **one-shot action** that forces:

* **all Wazuh agents**, or
* a specific **list of agents**

to reconnect to the manager using the `/agents/reconnect` API endpoint.

Use it when you want to:

* force agents to re-establish their connection (e.g. after network, config or certificate changes),
* optionally target only selected agents,
* store the API response (message, counts, error code, timestamp) in Terraform as outputs.

> This resource does **not** represent persistent state.
> Each `apply` that creates this resource sends **a new “force reconnect” command**.

---

## Example Usage

### Force Reconnect All Agents

```hcl
resource "wazuh_agent_reconnect" "all_agents" {}

output "reconnect_all_message" {
  value = wazuh_agent_reconnect.all_agents.message
}
```

### Force Reconnect Specific Agents

```hcl
resource "wazuh_agent_reconnect" "some_agents" {
  agents_list = ["001", "002", "003"]
}

output "reconnect_some_message" {
  value = wazuh_agent_reconnect.some_agents.message
}

output "reconnect_some_total_affected" {
  value = wazuh_agent_reconnect.some_agents.total_affected
}

output "reconnect_some_error_code" {
  value = wazuh_agent_reconnect.some_agents.error_code
}
```

---

## Lifecycle & Behavior

### Create – Send Force Reconnect Command

On `terraform apply`, the resource:

* Builds the request URL:

  * If `agents_list` is **not set** or empty:

    ```http
    PUT /agents/reconnect
    ```

  * If `agents_list = ["001", "002"]`:

    ```http
    PUT /agents/reconnect?agents_list=001,002
    ```

* Sends a minimal JSON body (e.g. `{}`) for compatibility.

A typical successful response looks like:

```json
{
  "data": {
    "affected_items": [
      "001",
      "002",
      "003"
    ],
    "total_affected_items": 3,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "Force reconnect command was sent to all agents",
  "error": 0
}
```

The resource then:

* sets `message` from the response,
* sets `total_affected` and `total_failed`,
* sets `error_code` from the top-level `error`,
* sets `timestamp` to the current UTC time,
* sets `id` to a unique timestamp-based string (e.g. `20251108T210300Z`).

> ℹ️ A non-zero `error_code` does **not automatically fail** the Terraform apply – HTTP status is used as the success criterion. The `error_code` field is exposed so you can inspect or react to partial failures in your Terraform logic or external tooling.

---

### Read – No-op

On refresh/plan, `Read`:

* does **not** call Wazuh again,
* simply keeps the values stored at `Create` time.

This is intentional – a reconnect command is an instantaneous action, not something that can be “refreshed”.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** call any Wazuh API,
* only removes this action record from Terraform state.

The agents are **not** reconnected again or rolled back – the force reconnect has already happened.

---

## Arguments Reference

| Name          | Type         | Required    | Description                                                                                                                                                                                                 |
| ------------- | ------------ | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `agents_list` | list(string) | 🚫 optional | Optional list of agent IDs to force reconnect (e.g. `["001", "002"]`). If omitted, the Wazuh API will target **all agents**. Changes to this field force a new resource (and thus a new reconnect command). |

> 🔁 `agents_list` is marked as **ForceNew** – changing it will cause Terraform to destroy/recreate the resource, which triggers a new force reconnect with the new list of agents.

---

## Attributes Reference

| Name             | Type         | Description                                                                                 |
| ---------------- | ------------ | ------------------------------------------------------------------------------------------- |
| `id`             | string       | Unique ID of this reconnect action, generated from UTC timestamp (e.g. `20251108T210300Z`). |
| `agents_list`    | list(string) | The list of agents targeted by this reconnect action (if any).                              |
| `message`        | string       | Human-readable message returned by Wazuh after sending the force reconnect command.         |
| `total_affected` | int          | Number of agents for which the reconnect command was processed.                             |
| `total_failed`   | int          | Number of agents where the reconnect command failed.                                        |
| `error_code`     | int          | Raw `error` code from the Wazuh API response (`0` = success).                               |
| `timestamp`      | string       | UTC timestamp (RFC3339) when the reconnect request was sent via Terraform.                  |
