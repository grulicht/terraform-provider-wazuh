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

variable "wazuh_role_name" {
  type        = string
  description = "Name of the Wazuh role to be created or managed."
  default     = "readonly-role"
}

variable "wazuh_user_username" {
  type        = string
  description = "Username of the Wazuh user to be created or managed."
  default     = "test-user"
}

variable "wazuh_user_password" {
  type        = string
  description = "Password for the Wazuh user. Store securely (e.g., in Terraform Cloud variables or a secrets manager)."
  sensitive   = true
  default     = "Pa$$word123"
}

variable "wazuh_role_user_position" {
  type        = number
  description = "Position (priority) of the role assignment for the given Wazuh user. Typically 0."
  default     = 0
}

variable "wazuh_policy_role_position" {
  type    = number
  default = 0
}

variable "wazuh_policy_name" {
  type        = string
  description = "Name of the Wazuh security policy to be created or managed."
  default     = "normal_policy_from_tf"
}

variable "wazuh_policy_actions" {
  type        = list(string)
  description = "List of actions allowed or denied by this Wazuh policy."
  default     = []
}

variable "wazuh_policy_resources" {
  type        = list(string)
  description = "List of resources this Wazuh policy applies to."
  default     = []
}

variable "wazuh_policy_effect" {
  type        = string
  description = "Effect of the Wazuh policy. Typically 'allow' or 'deny'."
  default     = "allow"
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
