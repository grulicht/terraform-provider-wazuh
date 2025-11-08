output "node_restart_message" {
  value = wazuh_node_restart.all_nodes.message
}

output "node_restart_total_affected" {
  value = wazuh_node_restart.all_nodes.total_affected
}
