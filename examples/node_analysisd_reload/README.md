<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_node_analysisd_reload.all_nodes](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/node_analysisd_reload) | resource |
| [wazuh_node_analysisd_reload.some_nodes](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/node_analysisd_reload) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_nodes_list"></a> [wazuh\_nodes\_list](#input\_wazuh\_nodes\_list) | List of Wazuh node IDs to which this configuration or resource should apply. | `list(string)` | <pre>[<br/>  "node01",<br/>  "node02"<br/>]</pre> | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_analysisd_reload_message"></a> [analysisd\_reload\_message](#output\_analysisd\_reload\_message) | n/a |
| <a name="output_analysisd_reload_total_affected"></a> [analysisd\_reload\_total\_affected](#output\_analysisd\_reload\_total\_affected) | n/a |
<!-- END_TF_DOCS -->