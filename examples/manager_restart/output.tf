output "manager_restart_message" {
  value = wazuh_manager_restart.restart.message
}

output "manager_restart_affected" {
  value = wazuh_manager_restart.restart.total_affected
}
