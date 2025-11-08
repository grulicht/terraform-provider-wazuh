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

variable "wazuh_node_configuration_node_id" {
  type        = string
  description = "ID of the Wazuh node to which this configuration will be applied."
  default     = "node01"
}

variable "wazuh_node_configuration_xml" {
  type        = string
  description = "XML configuration to be applied to the specified Wazuh node."
  default     = <<EOT
<ossec_config>
  <global>
    <email_notification>no</email_notification>
  </global>
</ossec_config>
EOT
}
