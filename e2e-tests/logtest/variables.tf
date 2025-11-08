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

variable "wazuh_logtest_log_format" {
  type        = string
  description = "Log format used by the Wazuh logtest resource (e.g. syslog, json, etc.)."
  default     = "syslog"
}

variable "wazuh_logtest_location" {
  type        = string
  description = "Path to the log file used for Wazuh log testing."
  default     = "/var/log/app.log"
}

variable "wazuh_logtest_event" {
  type        = string
  description = "Sample log event to be sent to the Wazuh logtest API for rule/decoder testing."
  default     = "2025-11-08T12:00:00Z my-app ERROR something broke"
}
