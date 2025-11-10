<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_policy.normal_policy](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/policy) | resource |
| [wazuh_policy_role.readonly_policy](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/policy_role) | resource |
| [wazuh_role.readonly](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/role) | resource |
| [wazuh_role_user.user_roles](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/role_user) | resource |
| [wazuh_security_rule.normal_rule](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/security_rule) | resource |
| [wazuh_security_rule_role.readonly_security_rule](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/security_rule_role) | resource |
| [wazuh_user.user](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/user) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_policy_actions"></a> [wazuh\_policy\_actions](#input\_wazuh\_policy\_actions) | List of actions allowed or denied by this Wazuh policy. | `list(string)` | `[]` | no |
| <a name="input_wazuh_policy_effect"></a> [wazuh\_policy\_effect](#input\_wazuh\_policy\_effect) | Effect of the Wazuh policy. Typically 'allow' or 'deny'. | `string` | `"allow"` | no |
| <a name="input_wazuh_policy_name"></a> [wazuh\_policy\_name](#input\_wazuh\_policy\_name) | Name of the Wazuh security policy to be created or managed. | `string` | `"normal_policy_from_tf"` | no |
| <a name="input_wazuh_policy_resources"></a> [wazuh\_policy\_resources](#input\_wazuh\_policy\_resources) | List of resources this Wazuh policy applies to. | `list(string)` | `[]` | no |
| <a name="input_wazuh_policy_role_position"></a> [wazuh\_policy\_role\_position](#input\_wazuh\_policy\_role\_position) | n/a | `number` | `0` | no |
| <a name="input_wazuh_role_name"></a> [wazuh\_role\_name](#input\_wazuh\_role\_name) | Name of the Wazuh role to be created or managed. | `string` | `"readonly-role"` | no |
| <a name="input_wazuh_role_user_position"></a> [wazuh\_role\_user\_position](#input\_wazuh\_role\_user\_position) | Position (priority) of the role assignment for the given Wazuh user. Typically 0. | `number` | `0` | no |
| <a name="input_wazuh_security_rule_definition"></a> [wazuh\_security\_rule\_definition](#input\_wazuh\_security\_rule\_definition) | Definition value used inside the MATCH block of the Wazuh security rule. | `string` | `"normalRul"` | no |
| <a name="input_wazuh_security_rule_name"></a> [wazuh\_security\_rule\_name](#input\_wazuh\_security\_rule\_name) | Name of the Wazuh security rule to be created. | `string` | `"normalRuleFromTF"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
| <a name="input_wazuh_user_password"></a> [wazuh\_user\_password](#input\_wazuh\_user\_password) | Password for the Wazuh user. Store securely (e.g., in Terraform Cloud variables or a secrets manager). | `string` | `"Pa$$word123"` | no |
| <a name="input_wazuh_user_username"></a> [wazuh\_user\_username](#input\_wazuh\_user\_username) | Username of the Wazuh user to be created or managed. | `string` | `"test-user"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_security_rule_id"></a> [security\_rule\_id](#output\_security\_rule\_id) | n/a |
| <a name="output_wazuh_policy_id"></a> [wazuh\_policy\_id](#output\_wazuh\_policy\_id) | n/a |
| <a name="output_wazuh_readonly_role_id"></a> [wazuh\_readonly\_role\_id](#output\_wazuh\_readonly\_role\_id) | n/a |
| <a name="output_wazuh_role_user_message"></a> [wazuh\_role\_user\_message](#output\_wazuh\_role\_user\_message) | n/a |
| <a name="output_wazuh_user_id"></a> [wazuh\_user\_id](#output\_wazuh\_user\_id) | n/a |
<!-- END_TF_DOCS -->