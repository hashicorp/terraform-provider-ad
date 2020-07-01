variable name { default = "TestOU" }
variable path { default = "dc=yourdomain,dc=com" }
variable description { default = "some description" }
variable protected { default = false }

variable name { default = "test group" }
variable sam_account_name { default = "TESTGROUP" }
variable scope { default = "global" }
variable category { default = "security" }

resource "ad_ou" "o" {
  name        = var.name
  path        = var.path
  description = var.description
  protected   = var.protected
}


resource "ad_group" "g" {
  name             = var.name
  sam_account_name = var.sam_account_name
  scope            = var.scope
  category         = var.category
  container        = ad_ou.o.dn
}
