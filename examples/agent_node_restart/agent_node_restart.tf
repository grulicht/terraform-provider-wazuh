resource "wazuh_agent_node_restart" "node01_agents" {
  node_id = var.wazuh_agent_node_restart_node_id
}