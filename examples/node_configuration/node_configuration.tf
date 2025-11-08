resource "wazuh_node_configuration" "node01" {
  node_id           = var.wazuh_node_configuration_node_id
  configuration_xml = var.wazuh_node_configuration_xml
}
