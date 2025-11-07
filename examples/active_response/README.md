<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_active_response.restart_agent](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/active_response) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_active_response_agents_list"></a> [wazuh\_active\_response\_agents\_list](#input\_wazuh\_active\_response\_agents\_list) | List of Wazuh agent IDs to which the active response will be applied. | `list(string)` | <pre>[<br/>  "001"<br/>]</pre> | no |
| <a name="input_wazuh_active_response_arguments"></a> [wazuh\_active\_response\_arguments](#input\_wazuh\_active\_response\_arguments) | Arguments to be passed to the active response command. | `list(string)` | <pre>[<br/>  "now"<br/>]</pre> | no |
| <a name="input_wazuh_active_response_command"></a> [wazuh\_active\_response\_command](#input\_wazuh\_active\_response\_command) | The command to be executed as part of the Wazuh active response. Use exclamation mark (!) for built-in commands. | `string` | `"!restart.sh"` | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_active_response_message"></a> [active\_response\_message](#output\_active\_response\_message) | Message returned by the Wazuh API after running the Active Response command. |
| <a name="output_active_response_timestamp"></a> [active\_response\_timestamp](#output\_active\_response\_timestamp) | Timestamp when the Active Response command was executed. |
| <a name="output_active_response_total_affected"></a> [active\_response\_total\_affected](#output\_active\_response\_total\_affected) | Total number of agents affected by the command. |
<!-- END_TF_DOCS -->