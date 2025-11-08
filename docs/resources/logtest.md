# 🐞 **Resource Documentation: `wazuh_logtest`**

# wazuh_logtest

The `wazuh_logtest` resource allows you to run the **Wazuh logtest tool** via the `/logtest` API endpoint.

It lets you send a sample log/event to Wazuh and see:

* whether it raises an alert,
* which rules/decoders matched,
* diagnostic output & messages.

This resource represents a **one-time test action**.
On `apply`, it runs `logtest`; on `destroy`, it ends the logtest session via `/logtest/sessions/{token}`.

---

## Example Usage

### Basic Logtest Run

```hcl
resource "wazuh_logtest" "example" {
  log_format = "syslog"
  location   = "/var/log/app.log"
  event      = "2025-11-08T12:00:00Z my-app ERROR something broke"
}

output "logtest_alert" {
  value = wazuh_logtest.example.alert
}

output "logtest_messages" {
  value = wazuh_logtest.example.messages
}

output "logtest_output_raw" {
  value = wazuh_logtest.example.output
}
```

After `terraform apply`, you’ll see whether the event triggered an alert and the structured output from `logtest`.

---

## Supported `log_format` Values

According to the Wazuh API, the `log_format` field supports (among others):

* `syslog`
* `json`
* `snort-full`
* `squid`
* `eventlog`
* `eventchannel`
* `audit`
* `mysql_log`
* `postgresql_log`
* `nmapg`
* `iis`
* `command`
* `full_command`
* `djb-multilog`
* `multi-line`

You must provide one of these values (or any valid format accepted by your Wazuh version).

---

## Lifecycle & Behavior

* Creating a `wazuh_logtest` resource sends:

  ```http
  PUT /logtest
  Content-Type: application/json
  ```

  with a body like:

  ```json
  {
    "log_format": "syslog",
    "location": "/var/log/app.log",
    "event": "2025-11-08T12:00:00Z my-app ERROR something broke"
  }
  ```

  Wazuh responds with something like:

  ```json
  {
    "error": 0,
    "data": {
      "messages": [],
      "token": "8cd2d6d1",
      "output": {},
      "alert": false,
      "codemsg": 1
    }
  }
  ```

* The resource stores:

  * the **session token**,
  * list of **messages**,
  * whether there was an **alert**,
  * `codemsg` and raw `output`.

* Destroying the resource calls:

  ```http
  DELETE /logtest/sessions/{token}
  ```

  which ends the saved logtest session on the Wazuh side.

> ⚠️ Each `terraform apply` re-runs logtest with the specified event. Treat it as an **action**, not a persistent object.

---

## Arguments Reference

| Name         | Type   | Required               | Description                                                                                              |
| ------------ | ------ | ---------------------- | -------------------------------------------------------------------------------------------------------- |
| `log_format` | string | ✅ **Yes**              | Log format used by logtest (e.g. `syslog`, `json`, `eventchannel`, `command`, etc.).                     |
| `location`   | string | ✅ **Yes**              | Location/path string associated with the event (e.g. `/var/log/app.log`).                                |
| `event`      | string | ✅ **Yes**              | The log line / event content to evaluate via logtest.                                                    |
| `token`      | string | 🚫 optional / computed | Optional existing logtest session token. If omitted, Wazuh will create a new session and return a token. |

---

## Attributes Reference

| Name        | Type         | Description                                                                                  |
| ----------- | ------------ | -------------------------------------------------------------------------------------------- |
| `id`        | string       | Resource ID. If Wazuh returns a token, it is used as the ID; otherwise a timestamp-based ID. |
| `token`     | string       | Logtest session token returned by Wazuh.                                                     |
| `alert`     | bool         | `true` if the provided event raised an alert; `false` otherwise.                             |
| `codemsg`   | int          | Numeric code returned by logtest (e.g. success / failure code).                              |
| `messages`  | list(string) | Diagnostic or informational messages returned by logtest.                                    |
| `output`    | string       | Raw JSON `output` field from logtest serialized as a string.                                 |
| `timestamp` | string       | UTC timestamp when logtest was executed.                                                     |
