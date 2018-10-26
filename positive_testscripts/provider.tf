provider "opentelekomcloud" {
  user_name   = "${var.username}"
  password    = "${var.password}"
//  tenant_name = "${var.tenant_name}"
  tenant_id = "${var.tenant_id}"
  domain_name = "${var.domain_name}"
  auth_url    = "${var.endpoint}"
}
