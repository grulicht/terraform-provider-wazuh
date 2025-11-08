resource "wazuh_cdb_list" "cdb_list" {
  filename  = var.wazuh_cdb_list_filename
  overwrite = var.wazuh_cdb_list_overwrite
  content   = var.wazuh_cdb_list_content
}
