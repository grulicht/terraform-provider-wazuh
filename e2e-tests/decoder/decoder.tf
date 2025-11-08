resource "wazuh_decoder" "decoder" {
  filename  = var.wazuh_decoder_filename
  overwrite = var.wazuh_decoder_overwrite
  content   = var.wazuh_decoder_content
}