resource "wazuh_policy" "normal_policy" {
  name = var.wazuh_policy_name

  policy = jsonencode({
    actions   = var.wazuh_policy_actions
    resources = var.wazuh_policy_resources
    effect    = var.wazuh_policy_effect
  })
}

# 1) Role
resource "wazuh_role" "readonly" {
  name = var.wazuh_role_name
}

# 2) User
resource "wazuh_user" "user" {
  username = var.wazuh_user_username
  password = var.wazuh_user_password
}

# 3) User ↔ Role mapping
resource "wazuh_role_user" "user_roles" {
  user_id  = wazuh_user.user.user_id
  role_ids = [wazuh_role.readonly.role_id]
  position = var.wazuh_role_user_position
}

# 4) Role ↔ Policy mapping
resource "wazuh_policy_role" "readonly_policy" {
  role_id = wazuh_role.readonly.role_id

  policy_ids = [
    tonumber(wazuh_policy.normal_policy.policy_id)
  ]

  position = var.wazuh_policy_role_position
}

# 5) Security rule
resource "wazuh_security_rule" "normal_rule" {
  name = var.wazuh_security_rule_name

  rule = jsonencode({
    MATCH = {
      definition = var.wazuh_security_rule_definition
    }
  })
}

# 6) Role ↔ Security rule mapping
resource "wazuh_security_rule_role" "readonly_security_rule" {
  role_id = wazuh_role.readonly.role_id

  rule_ids = [
    tonumber(wazuh_security_rule.normal_rule.rule_id)
  ]
}
