<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent_node_restart.node01_agents](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_node_restart) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_node_restart_node_id"></a> [wazuh\_agent\_node\_restart\_node\_id](#input\_wazuh\_agent\_node\_restart\_node\_id) | ID of the Wazuh node whose agents should be restarted. | `string` | `"node01"` | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_node01_agents_restart_message"></a> [node01\_agents\_restart\_message](#output\_node01\_agents\_restart\_message) | n/a |
| <a name="output_node01_agents_restart_total_affected"></a> [node01\_agents\_restart\_total\_affected](#output\_node01\_agents\_restart\_total\_affected) | n/a |
<!-- END_TF_DOCS -->