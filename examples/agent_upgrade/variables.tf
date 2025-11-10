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

variable "wazuh_agent_upgrade_agents_list" {
  type        = list(string)
  description = "List of Wazuh agent IDs that should be upgraded."
  default     = ["001"]
}

variable "wazuh_agent_upgrade_version" {
  type        = string
  description = "Target Wazuh agent version to upgrade to."
  default     = "4.14.0"
}

variable "wazuh_agent_upgrade_wpk_repo" {
  type        = string
  description = "Optional WPK repository URL for the Wazuh agent upgrade."
  default     = "packages.wazuh.com/4.x"
}

variable "wazuh_agent_upgrade_use_http" {
  type        = bool
  description = "Whether to use HTTP instead of HTTPS for the WPK repository."
  default     = false
}

variable "wazuh_agent_upgrade_force" {
  type        = bool
  description = "Whether to force upgrade even if the agent is already on the target version."
  default     = true
}

variable "wazuh_agent_upgrade_package_type" {
  type        = string
  description = "Package type used for the Wazuh agent (e.g., deb, rpm)."
  default     = "deb"
}
