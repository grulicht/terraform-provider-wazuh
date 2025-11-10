resource "wazuh_user" "user" {
  username = var.wazuh_user_username
  password = var.wazuh_user_password
}
