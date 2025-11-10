resource "wazuh_security_config" "config" {
  auth_token_exp_timeout = var.wazuh_security_config_auth_token_exp_timeout
  rbac_mode              = var.wazuh_security_config_rbac_mode
}
