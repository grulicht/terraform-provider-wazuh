<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent.example](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_force_after_registration_time"></a> [wazuh\_agent\_force\_after\_registration\_time](#input\_wazuh\_agent\_force\_after\_registration\_time) | Time since registration (e.g. 2h) after which force insertion is allowed. | `string` | `"2h"` | no |
| <a name="input_wazuh_agent_force_disconnected_time_enabled"></a> [wazuh\_agent\_force\_disconnected\_time\_enabled](#input\_wazuh\_agent\_force\_disconnected\_time\_enabled) | Whether to enable force insertion based on agent disconnected time. | `bool` | `true` | no |
| <a name="input_wazuh_agent_force_disconnected_time_value"></a> [wazuh\_agent\_force\_disconnected\_time\_value](#input\_wazuh\_agent\_force\_disconnected\_time\_value) | Disconnected time threshold (e.g. 30m, 1h) after which force insertion is allowed. | `string` | `"30m"` | no |
| <a name="input_wazuh_agent_force_enabled"></a> [wazuh\_agent\_force\_enabled](#input\_wazuh\_agent\_force\_enabled) | Whether to enable force insertion behavior for the Wazuh agent. | `bool` | `true` | no |
| <a name="input_wazuh_agent_ip"></a> [wazuh\_agent\_ip](#input\_wazuh\_agent\_ip) | IP address or hostname of the Wazuh agent. | `string` | `"127.0.0.2"` | no |
| <a name="input_wazuh_agent_name"></a> [wazuh\_agent\_name](#input\_wazuh\_agent\_name) | Name of the Wazuh agent to be registered. | `string` | `"NewHost_"` | no |
| <a name="input_wazuh_agent_purge_on_destroy"></a> [wazuh\_agent\_purge\_on\_destroy](#input\_wazuh\_agent\_purge\_on\_destroy) | Whether to remove the agent from Wazuh when the Terraform resource is destroyed. | `bool` | `true` | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_agent_id"></a> [agent\_id](#output\_agent\_id) | n/a |
| <a name="output_agent_status"></a> [agent\_status](#output\_agent\_status) | n/a |
<!-- END_TF_DOCS -->