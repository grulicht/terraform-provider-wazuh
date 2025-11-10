output "wazuh_readonly_role_id" {
  value = wazuh_role.readonly.role_id
}

output "wazuh_user_id" {
  value = wazuh_user.user.user_id
}

output "wazuh_role_user_message" {
  value = wazuh_role_user.user_roles.message
}

output "wazuh_policy_id" {
  value = wazuh_policy.normal_policy.policy_id
}

output "security_rule_id" {
  value = wazuh_security_rule.normal_rule.rule_id
}
