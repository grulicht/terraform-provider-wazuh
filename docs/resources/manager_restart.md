# 🔄 **Resource Documentation: `wazuh_manager_restart`**

# wazuh_manager_restart

The `wazuh_manager_restart` resource models a **one-shot action** that restarts the Wazuh manager via the `/manager/restart` API endpoint.

Use it when you want to:

* restart the Wazuh manager as part of a Terraform run (e.g. after changing configuration, rules, decoders),
* and capture the response (message, counts, timestamp) in Terraform state/outputs.

> This resource does **not** represent long-lived state.
> Each `apply` triggers a **new restart request**.

---

## Example Usage

```hcl
resource "wazuh_manager_restart" "restart_now" {}

output "manager_restart_message" {
  value = wazuh_manager_restart.restart_now.message
}

output "manager_restart_total_affected" {
  value = wazuh_manager_restart.restart_now.total_affected
}

output "manager_restart_timestamp" {
  value = wazuh_manager_restart.restart_now.timestamp
}
```

* Running `terraform apply` will:

  * send a restart request to the Wazuh manager,
  * store the API response summary in state.

---

## Lifecycle & Behavior

### Create – Send Restart Request

On `terraform apply`, the resource:

* Calls:

  ```http
  PUT /manager/restart
  Content-Type: application/json
  Authorization: Bearer <token>
  ```

  with an empty JSON object `{}` as body (for compatibility).

* A successful response typically looks like:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 1,
      "failed_items": [],
      "total_failed_items": 0
    },
    "message": "Restart request sent to all specified nodes",
    "error": 0
  }
  ```

* The resource then:

  * sets `message` to the returned message,
  * sets `total_affected` and `total_failed`,
  * sets `timestamp` to the current UTC time,
  * sets `id` to a timestamp-based unique string (e.g. `20251108T201234Z`).

Each apply creates a **new “restart action” instance** in state.

---

### Read – No-op

On refresh/plan, the `Read` step:

* **does not call the API again**,
* just keeps the previously stored values (`message`, `total_affected`, `total_failed`, `timestamp`).

This is intentional – a restart is an instantaneous action, not a persistent resource to “refresh”.

---

### Delete – State Only

On `terraform destroy`, the resource:

* does **not** call any Wazuh API,
* simply removes the restart action from Terraform state.

The manager is **not** restarted again or “undone” – the action already happened at `Create` time.

---

## Arguments Reference

This resource is an action-only resource and does **not** accept any configuration arguments.

| Name     | Type | Required | Description                                                                                  |
| -------- | ---- | -------- | -------------------------------------------------------------------------------------------- |
| *(none)* | –    | –        | This resource currently has no input arguments. It always restarts the manager when created. |

---

## Attributes Reference

| Name             | Type   | Description                                                                                  |
| ---------------- | ------ | -------------------------------------------------------------------------------------------- |
| `id`             | string | Unique ID of this restart action. Generated from the UTC timestamp, e.g. `20251108T201234Z`. |
| `message`        | string | Human-readable message returned by Wazuh after sending the restart request.                  |
| `total_affected` | int    | Number of managers/nodes for which the restart request was successfully sent.                |
| `total_failed`   | int    | Number of managers/nodes where the restart request failed.                                   |
| `timestamp`      | string | UTC timestamp (RFC3339) when the restart request was triggered via Terraform.                |
