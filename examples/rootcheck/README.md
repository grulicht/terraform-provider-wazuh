<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_rootcheck.rootcheck](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/rootcheck) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_rootcheck_agent_id"></a> [wazuh\_rootcheck\_agent\_id](#input\_wazuh\_rootcheck\_agent\_id) | ID of the Wazuh agent on which rootcheck should be executed. | `string` | `"001"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_rootcheck_results_message"></a> [rootcheck\_results\_message](#output\_rootcheck\_results\_message) | n/a |
| <a name="output_rootcheck_results_total_affected"></a> [rootcheck\_results\_total\_affected](#output\_rootcheck\_results\_total\_affected) | n/a |
| <a name="output_rootcheck_scan_message"></a> [rootcheck\_scan\_message](#output\_rootcheck\_scan\_message) | n/a |
<!-- END_TF_DOCS -->