terraform {
  required_providers {
    wazuh = {
      source = "grulicht/wazuh"
    }
  }
}

provider "wazuh" {
  endpoint        = var.wazuh_endpoint
  user            = var.wazuh_api_user
  password        = var.wazuh_api_password
  skip_ssl_verify = var.wazuh_skip_ssl_verify
}
