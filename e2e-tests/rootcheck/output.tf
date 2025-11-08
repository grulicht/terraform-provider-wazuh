output "rootcheck_scan_message" {
  value = wazuh_rootcheck.rootcheck.scan_message
}

output "rootcheck_results_message" {
  value = wazuh_rootcheck.rootcheck.results_message
}

output "rootcheck_results_total_affected" {
  value = wazuh_rootcheck.rootcheck.results_total_affected
}
