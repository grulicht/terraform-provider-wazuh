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

variable "wazuh_cdb_list_filename" {
  type        = string
  description = "Filename (without extension) of the Wazuh CDB list."
  default     = "test_list"
}

variable "wazuh_cdb_list_overwrite" {
  type        = bool
  description = "Whether to overwrite the CDB list file if it already exists."
  default     = true
}

variable "wazuh_cdb_list_content" {
  type        = string
  description = "Content of the Wazuh CDB list, one entry per line in the format 'value:tag'."
  default     = <<EOT
192.168.1.10:allowed
192.168.1.11:denied
EOT
}
