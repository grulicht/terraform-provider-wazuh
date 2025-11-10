<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent_upgrade_custom.custom_wpk](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_upgrade_custom) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_upgrade_custom_agents_list"></a> [wazuh\_agent\_upgrade\_custom\_agents\_list](#input\_wazuh\_agent\_upgrade\_custom\_agents\_list) | List of Wazuh agent IDs that should be upgraded using a custom WPK file. | `list(string)` | <pre>[<br/>  "001"<br/>]</pre> | no |
| <a name="input_wazuh_agent_upgrade_custom_file_path"></a> [wazuh\_agent\_upgrade\_custom\_file\_path](#input\_wazuh\_agent\_upgrade\_custom\_file\_path) | Absolute path to the custom WPK file used for the Wazuh agent upgrade. | `string` | `"/var/ossec/updates/custom-agent-4.14.0.wpk"` | no |
| <a name="input_wazuh_agent_upgrade_custom_installer"></a> [wazuh\_agent\_upgrade\_custom\_installer](#input\_wazuh\_agent\_upgrade\_custom\_installer) | Optional custom installer script name to be used during the Wazuh agent upgrade. | `string` | `"upgrade.sh"` | no |
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_custom_upgrade_affected_items"></a> [custom\_upgrade\_affected\_items](#output\_custom\_upgrade\_affected\_items) | n/a |
| <a name="output_custom_upgrade_message"></a> [custom\_upgrade\_message](#output\_custom\_upgrade\_message) | n/a |
<!-- END_TF_DOCS -->