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

variable "wazuh_active_response_command" {
  type        = string
  description = "The command to be executed as part of the Wazuh active response. Use exclamation mark (!) for built-in commands."
  default     = "!restart.sh"
}

variable "wazuh_active_response_agents_list" {
  type        = list(string)
  description = "List of Wazuh agent IDs to which the active response will be applied."
  default     = ["001"]
}

variable "wazuh_active_response_arguments" {
  type        = list(string)
  description = "Arguments to be passed to the active response command."
  default     = ["now"]
}
