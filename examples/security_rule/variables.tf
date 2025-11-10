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

variable "wazuh_security_rule_name" {
  type        = string
  description = "Name of the Wazuh security rule to be created."
  default     = "normalRuleFromTF"
}

variable "wazuh_security_rule_definition" {
  type        = string
  description = "Definition value used inside the MATCH block of the Wazuh security rule."
  default     = "normalRul"
}
