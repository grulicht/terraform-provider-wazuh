<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent_restart.agents](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_restart) | resource |
| [wazuh_agent_restart.all_agents](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_restart) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_restart_agents_list"></a> [wazuh\_agent\_restart\_agents\_list](#input\_wazuh\_agent\_restart\_agents\_list) | List of Wazuh agent IDs that should be restarted. If empty, use the 'all\_agents' resource to restart all agents. | `list(string)` | <pre>[<br/>  "001"<br/>]</pre> | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_agent_restart_all_message"></a> [agent\_restart\_all\_message](#output\_agent\_restart\_all\_message) | n/a |
<!-- END_TF_DOCS -->