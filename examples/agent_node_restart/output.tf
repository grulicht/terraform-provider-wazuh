output "node01_agents_restart_message" {
  value = wazuh_agent_node_restart.node01_agents.message
}

output "node01_agents_restart_total_affected" {
  value = wazuh_agent_node_restart.node01_agents.total_affected
}
