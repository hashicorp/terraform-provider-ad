variable principal_name { default = "testuser" }
variable samaccountname { default = "testuser" }

resource "ad_user" "u" {
  principal_name   = var.principal_name
  sam_account_name = var.samaccountname
  display_name     = "Terraform Test User"
}

