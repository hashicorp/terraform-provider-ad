# basic example
variable principal_name { default = "testGmsa" }
variable samaccountname { default = "testGmsa" }

resource "ad_gmsa" "g" {
  name                                            = var.principal_name
  dns_host_name                                   = var.principal_name
}

# all Gmsa attributes
variable principal_name2 { default = "testGmsa2" }
variable samaccountname2 { default = "testGmsa2$" }
variable container       { default = "CN=Gmsas,DC=contoso,DC=com" }

resource "ad_gmsa" "g2" {
  name                                            = var.principal_name2
  sam_account_name                                = var.samaccountname2
  dns_host_name                                   = var.principal_name2  
  container                                       = var.container
  display_name                                    = "testGmsa2"
  description	                                  = "Some desc 2"
  delegated                                       = false
  managed_password_interval_in_days               = 15
  kerberos_encryption_type                        = [ "aes128","aes256" ]
  expiration                                      = "2021-12-30T00:00:00+00:00"
  service_principal_names                         = [
    "HTTP/Machine3.corp.contoso.com"
  ]
  principals_allowed_to_delegate_to_account       = [
    "CN=group1,DC=groups,DC=contoso,DC=com",
    "CN=computer1,DC=computers,DC=contoso,DC=com"
  ]
  principals_allowed_to_retrieve_managed_password = [
    "CN=group1,DC=groups,DC=contoso,DC=com",
    "CN=computer1,DC=computers,DC=contoso,DC=com"
  ]
}