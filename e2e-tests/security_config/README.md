<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_security_config.config](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/security_config) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_security_config_auth_token_exp_timeout"></a> [wazuh\_security\_config\_auth\_token\_exp\_timeout](#input\_wazuh\_security\_config\_auth\_token\_exp\_timeout) | Authentication token expiration timeout in seconds for the Wazuh security configuration. | `number` | `1800` | no |
| <a name="input_wazuh_security_config_rbac_mode"></a> [wazuh\_security\_config\_rbac\_mode](#input\_wazuh\_security\_config\_rbac\_mode) | RBAC mode for Wazuh security configuration. Typical values are 'white' or 'black'. | `string` | `"white"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_security_config_rbac_mode"></a> [security\_config\_rbac\_mode](#output\_security\_config\_rbac\_mode) | n/a |
| <a name="output_security_config_timeout"></a> [security\_config\_timeout](#output\_security\_config\_timeout) | n/a |
<!-- END_TF_DOCS -->