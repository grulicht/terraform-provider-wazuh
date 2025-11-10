resource "wazuh_agent_upgrade" "selected_agents" {
  agents_list     = var.wazuh_agent_upgrade_agents_list
  upgrade_version = var.wazuh_agent_upgrade_version

  # Optional parameters:
  # wpk_repo     = var.wazuh_agent_upgrade_wpk_repo
  # use_http     = var.wazuh_agent_upgrade_use_http
  # force        = var.wazuh_agent_upgrade_force
  # package_type = var.wazuh_agent_upgrade_package_type
}
