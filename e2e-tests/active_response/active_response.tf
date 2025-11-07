resource "wazuh_active_response" "restart_agent" {
  command     = var.wazuh_active_response_command
  agents_list = var.wazuh_active_response_agents_list
  arguments   = var.wazuh_active_response_arguments
}
