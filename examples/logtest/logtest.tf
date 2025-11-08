resource "wazuh_logtest" "logtest" {
  log_format = var.wazuh_logtest_log_format
  location   = var.wazuh_logtest_location
  event      = var.wazuh_logtest_event
}
