output "custom_upgrade_message" {
  value = wazuh_agent_upgrade.selected_agents.message
}

output "custom_upgrade_affected_items" {
  value = wazuh_agent_upgrade.selected_agents.affected_items
}
