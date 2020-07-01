variable "domain" { default = "yourdomain.com" }
variable "name" { default = "tfGPO" }

resource "ad_gpo" "gpo" {
  name   = var.name
  domain = var.domain
}
