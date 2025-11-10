# Restart all agents
resource "wazuh_agent_restart" "all_agents" {}

# Restart specific agents only
# resource "wazuh_agent_restart" "agents" {
#  agents_list = "001"
# }
