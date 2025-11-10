variable "wazuh_endpoint" {
  type        = string
  default     = "https://localhost:55000"
  description = "Wazuh URL"
}

variable "wazuh_api_user" {
  type        = string
  default     = "wazuh-wui"
  description = "Wazuh admin user"
  sensitive   = true
}

variable "wazuh_api_password" {
  type        = string
  default     = "MyS3cr37P450r.*-"
  description = "Wazuh password of admin user"
  sensitive   = true
}

variable "wazuh_skip_ssl_verify" {
  type        = bool
  description = "Whether to skip SSL certificate verification when connecting to the Wazuh API."
  default     = true
}

variable "wazuh_security_config_auth_token_exp_timeout" {
  type        = number
  description = "Authentication token expiration timeout in seconds for the Wazuh security configuration."
  default     = 1800
}

variable "wazuh_security_config_rbac_mode" {
  type        = string
  description = "RBAC mode for Wazuh security configuration. Typical values are 'white' or 'black'."
  default     = "white"
}
