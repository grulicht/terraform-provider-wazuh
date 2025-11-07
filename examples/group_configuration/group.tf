resource "wazuh_group" "group" {
  group_id = var.wazuh_group_id
}

resource "wazuh_group_configuration" "group_configuration" {
  group_id          = wazuh_group.group.group_id
  configuration_xml = var.wazuh_configuration_xml
}
