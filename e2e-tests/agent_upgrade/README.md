<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_agent_upgrade.selected_agents](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/agent_upgrade) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_agent_upgrade_agents_list"></a> [wazuh\_agent\_upgrade\_agents\_list](#input\_wazuh\_agent\_upgrade\_agents\_list) | List of Wazuh agent IDs that should be upgraded. | `list(string)` | <pre>[<br/>  "001"<br/>]</pre> | no |
| <a name="input_wazuh_agent_upgrade_force"></a> [wazuh\_agent\_upgrade\_force](#input\_wazuh\_agent\_upgrade\_force) | Whether to force upgrade even if the agent is already on the target version. | `bool` | `true` | no |
| <a name="input_wazuh_agent_upgrade_package_type"></a> [wazuh\_agent\_upgrade\_package\_type](#input\_wazuh\_agent\_upgrade\_package\_type) | Package type used for the Wazuh agent (e.g., deb, rpm). | `string` | `"deb"` | no |
| <a name="input_wazuh_agent_upgrade_use_http"></a> [wazuh\_agent\_upgrade\_use\_http](#input\_wazuh\_agent\_upgrade\_use\_http) | Whether to use HTTP instead of HTTPS for the WPK repository. | `bool` | `false` | no |
| <a name="input_wazuh_agent_upgrade_version"></a> [wazuh\_agent\_upgrade\_version](#input\_wazuh\_agent\_upgrade\_version) | Target Wazuh agent version to upgrade to. | `string` | `"4.14.0"` | no |
| <a name="input_wazuh_agent_upgrade_wpk_repo"></a> [wazuh\_agent\_upgrade\_wpk\_repo](#input\_wazuh\_agent\_upgrade\_wpk\_repo) | Optional WPK repository URL for the Wazuh agent upgrade. | `string` | `"packages.wazuh.com/4.x"` | no |
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