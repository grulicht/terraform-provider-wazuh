output "events_ingest_message" {
  value = wazuh_event.event.message
}

output "events_ingest_total_affected" {
  value = wazuh_event.event.total_affected
}