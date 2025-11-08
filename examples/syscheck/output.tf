output "syscheck_scan_message" {
  value = wazuh_syscheck.syscheck.scan_message
}

output "syscheck_results_message" {
  value = wazuh_syscheck.syscheck.results_message
}

output "syscheck_results_total_affected" {
  value = wazuh_syscheck.syscheck.results_total_affected
}
