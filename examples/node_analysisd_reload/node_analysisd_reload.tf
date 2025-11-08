resource "wazuh_node_analysisd_reload" "all_nodes" {
  # nodes_list omitted → all nodes
}

resource "wazuh_node_analysisd_reload" "some_nodes" {
  nodes_list = var.wazuh_nodes_list
}
