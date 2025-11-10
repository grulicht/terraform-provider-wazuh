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

variable "wazuh_agent_upgrade_custom_agents_list" {
  type        = list(string)
  description = "List of Wazuh agent IDs that should be upgraded using a custom WPK file."
  default     = ["001"]
}

variable "wazuh_agent_upgrade_custom_file_path" {
  type        = string
  description = "Absolute path to the custom WPK file used for the Wazuh agent upgrade."
  default     = "/var/ossec/updates/custom-agent-4.14.0.wpk"
}

variable "wazuh_agent_upgrade_custom_installer" {
  type        = string
  description = "Optional custom installer script name to be used during the Wazuh agent upgrade."
  default     = "upgrade.sh"
}
