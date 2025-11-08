resource "wazuh_rule" "rule" {
  filename  = var.wazuh_rule_filename
  overwrite = var.wazuh_rule_overwrite
  content   = var.wazuh_rule_content
}
