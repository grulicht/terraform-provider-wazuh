resource "wazuh_agent" "example" {
  name             = var.wazuh_agent_name
  ip               = var.wazuh_agent_ip
  purge_on_destroy = var.wazuh_agent_purge_on_destroy

  # Force insertion behaviour
  force_enabled                   = var.wazuh_agent_force_enabled
  force_disconnected_time_enabled = var.wazuh_agent_force_disconnected_time_enabled
  force_disconnected_time_value   = var.wazuh_agent_force_disconnected_time_value
  force_after_registration_time   = var.wazuh_agent_force_after_registration_time
}
