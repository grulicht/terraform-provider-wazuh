<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_decoder.decoder](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/decoder) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_decoder_content"></a> [wazuh\_decoder\_content](#input\_wazuh\_decoder\_content) | XML content of the Wazuh decoder definition. | `string` | `"<group name=\"local,syslog,\">\n  <decoder name=\"local-example\">\n    <prematch>my-app</prematch>\n  </decoder>\n</group>\n"` | no |
| <a name="input_wazuh_decoder_filename"></a> [wazuh\_decoder\_filename](#input\_wazuh\_decoder\_filename) | Filename of the local Wazuh decoder XML file. | `string` | `"local_decoder.xml"` | no |
| <a name="input_wazuh_decoder_overwrite"></a> [wazuh\_decoder\_overwrite](#input\_wazuh\_decoder\_overwrite) | Whether to overwrite the decoder file if it already exists. | `bool` | `true` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
<!-- END_TF_DOCS -->