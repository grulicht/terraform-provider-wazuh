resource "wazuh_agent_group" "example" {
  agent_id           = var.wazuh_agent_group_agent_id
  group_id           = var.wazuh_agent_group_group_id
  force_single_group = var.wazuh_agent_group_force_single_group
}
