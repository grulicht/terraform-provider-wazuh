output "active_response_message" {
  description = "Message returned by the Wazuh API after running the Active Response command."
  value       = wazuh_active_response.restart_agent.message
}

output "active_response_total_affected" {
  description = "Total number of agents affected by the command."
  value       = wazuh_active_response.restart_agent.total_affected
}

output "active_response_timestamp" {
  description = "Timestamp when the Active Response command was executed."
  value       = wazuh_active_response.restart_agent.timestamp
}
