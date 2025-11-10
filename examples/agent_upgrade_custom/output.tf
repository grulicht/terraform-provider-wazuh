output "custom_upgrade_message" {
  value = wazuh_agent_upgrade_custom.custom_wpk.message
}

output "custom_upgrade_affected_items" {
  value = wazuh_agent_upgrade_custom.custom_wpk.affected_items
}
