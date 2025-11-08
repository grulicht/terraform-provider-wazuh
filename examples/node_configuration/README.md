<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_node_configuration.node01](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/node_configuration) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_node_configuration_node_id"></a> [wazuh\_node\_configuration\_node\_id](#input\_wazuh\_node\_configuration\_node\_id) | ID of the Wazuh node to which this configuration will be applied. | `string` | `"node01"` | no |
| <a name="input_wazuh_node_configuration_xml"></a> [wazuh\_node\_configuration\_xml](#input\_wazuh\_node\_configuration\_xml) | XML configuration to be applied to the specified Wazuh node. | `string` | `"<ossec_config>\n  <global>\n    <email_notification>no</email_notification>\n  </global>\n</ossec_config>\n"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
<!-- END_TF_DOCS -->