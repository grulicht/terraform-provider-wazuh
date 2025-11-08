
resource "wazuh_event" "event" {
  events = var.wazuh_event_events
}
