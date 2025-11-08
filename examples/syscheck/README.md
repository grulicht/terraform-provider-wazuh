<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_syscheck.syscheck](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/syscheck) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
| <a name="input_wazuh_syscheck_agent_id"></a> [wazuh\_syscheck\_agent\_id](#input\_wazuh\_syscheck\_agent\_id) | ID of the Wazuh agent on which syscheck should be executed. | `string` | `"001"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_syscheck_results_message"></a> [syscheck\_results\_message](#output\_syscheck\_results\_message) | n/a |
| <a name="output_syscheck_results_total_affected"></a> [syscheck\_results\_total\_affected](#output\_syscheck\_results\_total\_affected) | n/a |
| <a name="output_syscheck_scan_message"></a> [syscheck\_scan\_message](#output\_syscheck\_scan\_message) | n/a |
<!-- END_TF_DOCS -->