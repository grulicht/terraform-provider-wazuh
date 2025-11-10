<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_security_rule.normal_rule](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/security_rule) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_security_rule_definition"></a> [wazuh\_security\_rule\_definition](#input\_wazuh\_security\_rule\_definition) | Definition value used inside the MATCH block of the Wazuh security rule. | `string` | `"normalRul"` | no |
| <a name="input_wazuh_security_rule_name"></a> [wazuh\_security\_rule\_name](#input\_wazuh\_security\_rule\_name) | Name of the Wazuh security rule to be created. | `string` | `"normalRuleFromTF"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_security_rule_id"></a> [security\_rule\_id](#output\_security\_rule\_id) | n/a |
<!-- END_TF_DOCS -->