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

variable "wazuh_decoder_filename" {
  type        = string
  description = "Filename of the local Wazuh decoder XML file."
  default     = "local_decoder.xml"
}

variable "wazuh_decoder_overwrite" {
  type        = bool
  description = "Whether to overwrite the decoder file if it already exists."
  default     = true
}

variable "wazuh_decoder_content" {
  type        = string
  description = "XML content of the Wazuh decoder definition."
  default     = <<EOT
<group name="local,syslog,">
  <decoder name="local-example">
    <prematch>my-app</prematch>
  </decoder>
</group>
EOT
}
