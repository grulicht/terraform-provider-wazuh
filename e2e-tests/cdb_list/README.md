<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_cdb_list.cdb_list](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/cdb_list) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_cdb_list_content"></a> [wazuh\_cdb\_list\_content](#input\_wazuh\_cdb\_list\_content) | Content of the Wazuh CDB list, one entry per line in the format 'value:tag'. | `string` | `"192.168.1.10:allowed\n192.168.1.11:denied\n"` | no |
| <a name="input_wazuh_cdb_list_filename"></a> [wazuh\_cdb\_list\_filename](#input\_wazuh\_cdb\_list\_filename) | Filename (without extension) of the Wazuh CDB list. | `string` | `"test_list"` | no |
| <a name="input_wazuh_cdb_list_overwrite"></a> [wazuh\_cdb\_list\_overwrite](#input\_wazuh\_cdb\_list\_overwrite) | Whether to overwrite the CDB list file if it already exists. | `bool` | `true` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
<!-- END_TF_DOCS -->