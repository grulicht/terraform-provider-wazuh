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

variable "wazuh_agent_name" {
  type        = string
  description = "Name of the Wazuh agent to be registered."
  default     = "NewHost_"
}

variable "wazuh_agent_ip" {
  type        = string
  description = "IP address or hostname of the Wazuh agent."
  default     = "127.0.0.2"
}

variable "wazuh_agent_purge_on_destroy" {
  type        = bool
  description = "Whether to remove the agent from Wazuh when the Terraform resource is destroyed."
  default     = true
}


variable "wazuh_agent_force_enabled" {
  type        = bool
  description = "Whether to enable force insertion behavior for the Wazuh agent."
  default     = true
}

variable "wazuh_agent_force_disconnected_time_enabled" {
  type        = bool
  description = "Whether to enable force insertion based on agent disconnected time."
  default     = true
}

variable "wazuh_agent_force_disconnected_time_value" {
  type        = string
  description = "Disconnected time threshold (e.g. 30m, 1h) after which force insertion is allowed."
  default     = "30m"
}

variable "wazuh_agent_force_after_registration_time" {
  type        = string
  description = "Time since registration (e.g. 2h) after which force insertion is allowed."
  default     = "2h"
}
