<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_wazuh"></a> [wazuh](#provider\_wazuh) | n/a |

## Resources

| Name | Type |
|------|------|
| [wazuh_event.event](https://registry.terraform.io/providers/grulicht/wazuh/latest/docs/resources/event) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_wazuh_api_password"></a> [wazuh\_api\_password](#input\_wazuh\_api\_password) | Wazuh password of admin user | `string` | `"MyS3cr37P450r.*-"` | no |
| <a name="input_wazuh_api_user"></a> [wazuh\_api\_user](#input\_wazuh\_api\_user) | Wazuh admin user | `string` | `"wazuh-wui"` | no |
| <a name="input_wazuh_endpoint"></a> [wazuh\_endpoint](#input\_wazuh\_endpoint) | Wazuh URL | `string` | `"https://localhost:55000"` | no |
| <a name="input_wazuh_event_events"></a> [wazuh\_event\_events](#input\_wazuh\_event\_events) | List of event payloads to be sent as Wazuh events. Values can be plain strings or JSON-encoded strings. | `list(string)` | <pre>[<br/>  "Event value 1",<br/>  "{\"someKey\": \"Event value 2\"}"<br/>]</pre> | no |
| <a name="input_wazuh_skip_ssl_verify"></a> [wazuh\_skip\_ssl\_verify](#input\_wazuh\_skip\_ssl\_verify) | Whether to skip SSL certificate verification when connecting to the Wazuh API. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_events_ingest_message"></a> [events\_ingest\_message](#output\_events\_ingest\_message) | n/a |
| <a name="output_events_ingest_total_affected"></a> [events\_ingest\_total\_affected](#output\_events\_ingest\_total\_affected) | n/a |
<!-- END_TF_DOCS -->