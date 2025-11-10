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

variable "wazuh_agent_group_agent_id" {
  type        = string
  description = "ID of the Wazuh agent that should be assigned to the group."
  default     = "001"
}

variable "wazuh_agent_group_group_id" {
  type        = string
  description = "ID of the Wazuh group to which the agent should be assigned."
  default     = "NewGroup_1"
}

variable "wazuh_agent_group_force_single_group" {
  type        = bool
  description = "Whether to force the agent to belong only to this single group (removing it from other groups)."
  default     = true
}
