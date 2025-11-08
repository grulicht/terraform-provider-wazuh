<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_rule.rule](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/rule) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_rule_content"></a> [wazuh\_rule\_content](#input\_wazuh\_rule\_content) | XML content of the Wazuh rule group and rules to be deployed. | `string` | `"<group name=\"local,\">\n  <rule id=\"100100\" level=\"3\">\n    <if_sid>100000</if_sid>\n    <match>my-custom-rule</match>\n    <description>My custom local rule</description>\n  </rule>\n</group>\n"` | no |
| <a name="input_wazuh_rule_filename"></a> [wazuh\_rule\_filename](#input\_wazuh\_rule\_filename) | Filename of the local Wazuh rule XML file. | `string` | `"local_rules.xml"` | no |
| <a name="input_wazuh_rule_overwrite"></a> [wazuh\_rule\_overwrite](#input\_wazuh\_rule\_overwrite) | Whether to overwrite the rule file if it already exists. | `bool` | `true` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |
<!-- END_TF_DOCS -->