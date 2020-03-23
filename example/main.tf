provider "msad" {
    bind_username = "kyriakos@yourdomain.com"
    bind_password = "xYbb6qyQUsrUE82qo0DubA=="
    dc_hostname = "127.0.0.1"
    dc_port = 1636
    proto = "ldaps"
    allow_insecure_certs = "true"
}

resource "msad_user" "a" {
    domain_dn = "dc=yourdomain,dc=com"
    display_name = "Test User"
    principal_name = "testuser@yourdomain.com"
    sam_account_name = "testuser"
    initial_password = "test12345!!!"
    change_at_next_login = false
}

#data "msad_user" "u" {
#    domain_dn = "Dc=yourdomain,Dc=com"
#    user_dn = "CN=Test User,CN=Users,Dc=yourdomain,Dc=com"
#}

data "msad_domain" "d" {
    netbios_name = "YOURDOMAIN"
}

#output "blah" {
#  value = data.msad_user.u
#}

output "domain" {
    value = data.msad_domain.d
}

