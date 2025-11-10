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
