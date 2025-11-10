resource "wazuh_policy" "normal_policy" {
  name = var.wazuh_policy_name

  policy = jsonencode({
    actions   = var.wazuh_policy_actions
    resources = var.wazuh_policy_resources
    effect    = var.wazuh_policy_effect
  })
}
