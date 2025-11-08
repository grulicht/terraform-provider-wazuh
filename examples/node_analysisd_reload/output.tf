output "analysisd_reload_message" {
  value = wazuh_node_analysisd_reload.all_nodes.message
}

output "analysisd_reload_total_affected" {
  value = wazuh_node_analysisd_reload.all_nodes.total_affected
}
