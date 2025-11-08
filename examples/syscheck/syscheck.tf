resource "wazuh_syscheck" "syscheck" {
  agent_id = var.wazuh_syscheck_agent_id
}
