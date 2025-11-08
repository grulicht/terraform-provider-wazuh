# 🌐 **Resource Documentation: `wazuh_event`**

# wazuh_event

The `wazuh_event` resource allows you to **ingest security events** into Wazuh’s analysisd via the **`/events`** API endpoint.

This resource is useful when you want to push custom events, logs, or JSON payloads into Wazuh for correlation, alerting, or testing rules and decoders – all directly from Terraform.

---

## Example Usage

### Ingest a Small Batch of Events

```hcl
resource "wazuh_event" "sample" {
  events = [
    "Event value 1",
    "{\"someKey\": \"Event value 2\"}",
  ]
}
```

### Ingest JSON-Formatted Security Events

```hcl
resource "wazuh_event" "login_events" {
  events = [
    "{\"username\": \"alice\", \"action\": \"login_failed\", \"source_ip\": \"10.0.0.10\"}",
    "{\"username\": \"bob\", \"action\": \"login_success\", \"source_ip\": \"10.0.0.11\"}",
  ]
}
```

### Outputs (optional)

```hcl
output "wazuh_events_message" {
  value = wazuh_event.sample.message
}

output "wazuh_events_total_affected" {
  value = wazuh_event.sample.total_affected
}

output "wazuh_events_total_failed" {
  value = wazuh_event.sample.total_failed
}
```

---

## Lifecycle & Behavior

* Creating a `wazuh_event` resource triggers a:

  ```http
  POST /events
  Content-Type: application/json
  ```

  with a request body like:

  ```json
  {
    "events": [
      "Event value 1",
      "{\"someKey\": \"Event value 2\"}"
    ]
  }
  ```

* Wazuh forwards the events to **analysisd**, and returns a response such as:

  ```json
  {
    "data": {
      "affected_items": [],
      "total_affected_items": 2,
      "total_failed_items": 0,
      "failed_items": []
    },
    "message": "All events were forwarded to analisysd",
    "error": 0
  }
  ```

* `terraform destroy` **does not delete events** from Wazuh – it only removes the resource from Terraform state.

> ⚠️ **Note:** This resource represents a **one-time ingestion** action. Each apply sends a new batch of events.

---

## Arguments Reference

| Name     | Type         | Required  | Description                                                                                   |
| -------- | ------------ | --------- | --------------------------------------------------------------------------------------------- |
| `events` | list(string) | ✅ **Yes** | List of events to ingest. Each element is a string (plain text or JSON). Max 100 per request. |

---

## Attributes Reference

| Name             | Type   | Description                                                                |
| ---------------- | ------ | -------------------------------------------------------------------------- |
| `id`             | string | Internal Terraform resource ID (based on timestamp).                       |
| `message`        | string | Human-readable description from Wazuh (e.g. “All events were forwarded…”). |
| `total_affected` | int    | Number of events successfully processed.                                   |
| `total_failed`   | int    | Number of events that failed to process.                                   |
| `timestamp`      | string | UTC timestamp when the events were ingested.                               |
