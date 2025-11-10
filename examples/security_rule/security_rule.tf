resource "wazuh_security_rule" "normal_rule" {
  name = var.wazuh_security_rule_name

  rule = jsonencode({
    MATCH = {
      definition = var.wazuh_security_rule_definition
    }
  })
}
