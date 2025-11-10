resource "wazuh_agent_upgrade_custom" "custom_wpk" {
  agents_list = var.wazuh_agent_upgrade_custom_agents_list
  file_path   = var.wazuh_agent_upgrade_custom_file_path

  # optional
  # installer   = var.wazuh_agent_upgrade_custom_installer
}
