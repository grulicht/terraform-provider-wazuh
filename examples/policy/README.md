<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_policy.normal_policy](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/policy) | resource |

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
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_wazuh_policy_id"></a> [wazuh\_policy\_id](#output\_wazuh\_policy\_id) | n/a |
<!-- END_TF_DOCS -->