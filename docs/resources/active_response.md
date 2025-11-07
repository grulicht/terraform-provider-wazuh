# 🧠 **Resource Documentation: `wazuh_active_response`**

# wazuh_active_response

The `wazuh_active_response` resource allows you to **run Active Response commands** on one or more Wazuh agents via the **Wazuh Manager REST API**.
This lets you execute remote commands, scripts, or remediation actions (e.g., restart agents, block IPs, remove files, etc.) directly from Terraform.

---

## Example Usage

### Run a Restart Command on a Single Agent

```hcl
resource "wazuh_active_response" "restart_agent" {
  command      = "!restart.sh"
  agents_list  = ["001"]
  arguments    = ["now"]
}
```

### Run a Custom Script on Multiple Agents

```hcl
resource "wazuh_active_response" "firewall_block" {
  command      = "!firewall-drop.sh"
  agents_list  = ["001", "002"]
  arguments    = ["192.168.1.100"]
}
```

### View Outputs (optional)

You can capture execution details with an `output.tf`:

```hcl
output "active_response_message" {
  value = wazuh_active_response.restart_agent.message
}

output "active_response_total_affected" {
  value = wazuh_active_response.restart_agent.total_affected
}

output "active_response_timestamp" {
  value = wazuh_active_response.restart_agent.timestamp
}
```

---

## Lifecycle & Behavior

* Running a new `wazuh_active_response` resource triggers a `PUT /active-response` API call with payload like:

  ```json
  {
    "command": "!restart.sh",
    "arguments": ["now"]
  }
  ```

  and query parameter:

  ```
  ?agents_list=001,002
  ```

* The command is executed immediately on the specified agent(s).

* Terraform does not track execution results over time — each apply creates a new “execution record”.

* `terraform destroy` only removes the record from state (no API call is made).

---

## Arguments Reference

| Name          | Type         | Required    | Description                                                                                |
| ------------- | ------------ | ----------- | ------------------------------------------------------------------------------------------ |
| `command`     | string       | ✅ **Yes**   | Active Response command or script to execute (e.g., `!restart.sh`, `!firewall-drop.sh`).   |
| `arguments`   | list(string) | 🚫 optional | List of arguments to pass to the command.                                                  |
| `agents_list` | list(string) | 🚫 optional | List of agent IDs on which to run the command. If omitted, the command runs on all agents. |

---

## Attributes Reference

| Name             | Type   | Description                                                                                      |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------ |
| `id`             | string | Internal Terraform resource ID (generated automatically).                                        |
| `message`        | string | Message returned by Wazuh after command execution (e.g., `"AR command was sent to all agents"`). |
| `total_affected` | int    | Number of agents affected by the command.                                                        |
| `timestamp`      | string | UTC timestamp when the command was executed.                                                     |

---

### Destroy

```bash
terraform destroy
```

> 🗑️ **No API request is sent.**
> Terraform simply removes the resource from its state file — Active Response commands cannot be “undone” via API.

---

## Example Use Cases

| Scenario             | Command Example     | Description                                      |
| -------------------- | ------------------- | ------------------------------------------------ |
| Restart agent        | `!restart.sh`       | Restarts the Wazuh agent remotely.               |
| Drop IP via firewall | `!firewall-drop.sh` | Blocks an IP address using firewall rules.       |
| Remove file          | `!delete-file.sh`   | Deletes a specified file from the agent.         |
| Custom remediation   | `!custom-script.sh` | Executes your own custom Active Response script. |
