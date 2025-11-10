output "security_config_timeout" {
  value = wazuh_security_config.config.auth_token_exp_timeout
}

output "security_config_rbac_mode" {
  value = wazuh_security_config.config.rbac_mode
}
