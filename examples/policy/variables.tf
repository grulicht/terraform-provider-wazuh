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
