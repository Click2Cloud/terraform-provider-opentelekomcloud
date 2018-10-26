### OpenTelekomCloud Credentials
variable "username" {
  # If you don't fill this in, you will be prompted for it
  default = "lizhonghua"
}

variable "password" {
  # If you don't fill this in, you will be prompted for it
  default = "Tools@1234"
}

variable "domain_name" {
  # If you don't fill this in, you will be prompted for it
  default = "OTC00000000001000010501"
}

/*variable "tenant_name" {
  default = "eu-de_terraform"
}*/

variable "tenant_id" {
  default = "17fbda95add24720a4038ba4b1c705ed"
}

variable "endpoint" {
  default = "https://iam.eu-de.otc.t-systems.com/v3"
}

### OTC Specific Settings
variable "external_network" {
  default = "admin_external_net"
}

### Project Settings
variable "project" {
  default = "terraform"
}

variable "subnet_cidr" {
  default = "192.168.10.0/24"
}

variable "ssh_pub_key" {
  default = "~/.ssh/id_rsa.pub"
}

### DNS Settings
variable "dnszone" {
  default = ""
}

variable "dnsname" {
  default = "webserver"
}

### VM (Instance) Settings
variable "instance_count" {
  default = "1"
}

variable "disk_size_gb" {
  default = "0"
}

variable "flavor_name" {
  default = "s1.medium"
}

variable "image_name" {
  default = "Standard_CentOS_7_latest"
}

variable "endpoint_email" {
  default = "mailtest@gmail.com"
}

variable "endpoint_sms" {
  default = "+8613600000000"
}