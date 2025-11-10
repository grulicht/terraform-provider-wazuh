<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_user.user](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/user) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
| <a name="input_wazuh_user_password"></a> [wazuh\_user\_password](#input\_wazuh\_user\_password) | Password for the Wazuh user. Store securely (e.g., in Terraform Cloud variables or a secrets manager). | `string` | `"Pa$$word123"` | no |
| <a name="input_wazuh_user_username"></a> [wazuh\_user\_username](#input\_wazuh\_user\_username) | Username of the Wazuh user to be created or managed. | `string` | `"test-user"` | no |
<!-- END_TF_DOCS -->