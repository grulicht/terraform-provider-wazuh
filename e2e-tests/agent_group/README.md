<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent_group.example](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_group) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_group_agent_id"></a> [wazuh\_agent\_group\_agent\_id](#input\_wazuh\_agent\_group\_agent\_id) | ID of the Wazuh agent that should be assigned to the group. | `string` | `"001"` | no |
| <a name="input_wazuh_agent_group_force_single_group"></a> [wazuh\_agent\_group\_force\_single\_group](#input\_wazuh\_agent\_group\_force\_single\_group) | Whether to force the agent to belong only to this single group (removing it from other groups). | `bool` | `true` | no |
| <a name="input_wazuh_agent_group_group_id"></a> [wazuh\_agent\_group\_group\_id](#input\_wazuh\_agent\_group\_group\_id) | ID of the Wazuh group to which the agent should be assigned. | `string` | `"NewGroup_1"` | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_agent_group_message"></a> [agent\_group\_message](#output\_agent\_group\_message) | n/a |
<!-- END_TF_DOCS -->