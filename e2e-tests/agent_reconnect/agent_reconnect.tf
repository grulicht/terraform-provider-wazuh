# Force reconnect all agents
resource "wazuh_agent_reconnect" "all_agents" {}

# Force reconnect a specific agent
# resource "wazuh_agent_reconnect" "agents" {
#  agents_list = ["001"]
# }