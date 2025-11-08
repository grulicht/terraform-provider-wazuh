<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_logtest.logtest](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/logtest) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_logtest_event"></a> [wazuh\_logtest\_event](#input\_wazuh\_logtest\_event) | Sample log event to be sent to the Wazuh logtest API for rule/decoder testing. | `string` | `"2025-11-08T12:00:00Z my-app ERROR something broke"` | no |
| <a name="input_wazuh_logtest_location"></a> [wazuh\_logtest\_location](#input\_wazuh\_logtest\_location) | Path to the log file used for Wazuh log testing. | `string` | `"/var/log/app.log"` | no |
| <a name="input_wazuh_logtest_log_format"></a> [wazuh\_logtest\_log\_format](#input\_wazuh\_logtest\_log\_format) | Log format used by the Wazuh logtest resource (e.g. syslog, json, etc.). | `string` | `"syslog"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_logtest_alert"></a> [logtest\_alert](#output\_logtest\_alert) | n/a |
| <a name="output_logtest_messages"></a> [logtest\_messages](#output\_logtest\_messages) | n/a |
| <a name="output_logtest_output_raw"></a> [logtest\_output\_raw](#output\_logtest\_output\_raw) | n/a |
<!-- END_TF_DOCS -->