resource "wazuh_agent_restart_group" "agents" {
  group_id = var.wazuh_agent_restart_group_id
}
