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

variable "wazuh_rule_filename" {
  type        = string
  description = "Filename of the local Wazuh rule XML file."
  default     = "local_rules.xml"
}

variable "wazuh_rule_overwrite" {
  type        = bool
  description = "Whether to overwrite the rule file if it already exists."
  default     = true
}

variable "wazuh_rule_content" {
  type        = string
  description = "XML content of the Wazuh rule group and rules to be deployed."
  default     = <<EOT
<group name="local,">
  <rule id="100100" level="3">
    <if_sid>100000</if_sid>
    <match>my-custom-rule</match>
    <description>My custom local rule</description>
  </rule>
</group>
EOT
}
