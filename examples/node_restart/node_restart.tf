resource "wazuh_node_restart" "all_nodes" {
  # nodes_list = []  # if omitted or empty → restart all nodes
}

resource "wazuh_node_restart" "some_nodes" {
  nodes_list = var.wazuh_nodes_list
}
