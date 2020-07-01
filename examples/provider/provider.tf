resource "ad_ou" "o" { 
    name = "gplinktestOU"
    path = "dc=yourdomain,dc=com"
    description = "OU for gplink tests"
    protected = false
}
    
resource "ad_gpo" "g" {
    name        = "gplinktestGPO"
    domain      = "yourdomain.com"
    description = "gpo for gplink tests"
    status      = "AllSettingsEnabled"
}

resource "ad_gplink" "og" { 
    gpo_guid = ad_gpo.g.id
    target_dn = ad_ou.o.dn
    enforced = true
    enabled = true
    order = 0
}




#variable"domain"      { default = "yourdomain.com" }
#variable "name"        { default = "tfgpo" }
#
#resource "ad_gpo" "gpo" {
#    name        = var.name
#    domain      = var.domain
#}
#
#data "ad_gpo" "g" {
#    name = var.name
#}
#


#variable "name" { default = "test123"}
#variable "pre2kname" { default = "haha123" }
#
#resource "ad_computer" "c" {
#	name = var.name
#	pre2kname = var.pre2kname
#}
#
#	data "ad_computer" "dsc" {
#		guid = ad_computer.c.guid
#	}
#
# variable name { default = "hahaou" }
# variable path { default = "dc=yourdomain,dc=com" }
# variable description { default = "some description" }
# variable protected { default = false }
# 
# variable domain_dn { default = "dc=yourdomain,dc=com" }
# variable principal_name { default = "testuser" }
# variable password { default = "SuperSecurePassword123!!" }
# variable samaccountname { default = "testuser" }
# 
# variable display_name { default = "test group" }
# variable sam_account_name { default = "testgroup" }
# variable scope { default = "global" }
# variable type { default = "security" }
# 
# variable domain      { default = "yourdomain.com" }
# variable gpo_name        { default = "gaa123r" }
# 
# resource "ad_ou" "o1" { 
#     name = "another-${var.name}"
#     path = var.path
#     description = var.description
#     protected = var.protected
# }
# 
# resource "ad_ou" "o" { 
#     name = var.name
#     path = var.path
#     description = var.description
#     protected = var.protected
# }
# 
# 
# 	data "ad_ou" "ods" {
# 		dn = ad_ou.o.dn
# 	}
# 
# resource "ad_user" "u" {
# 	domain_dn = var.domain_dn
# 	principal_name = var.principal_name
# 	sam_account_name = var.samaccountname
# 	initial_password = var.password
# 	display_name = "Terraform Test User"
#     user_container = ad_ou.o1.dn
# }
# 
# resource "ad_group" "g" {
#     domain_dn = var.domain_dn
#     display_name = var.display_name
#     sam_account_name = var.sam_account_name
#     scope = var.scope
#     type = var.type
#     container = ad_ou.o.dn
# }
# 
# 
# 
# resource "ad_gpo" "gpo" {
#     name        = var.gpo_name
#     domain      = var.domain
# }
# 
# resource "ad_gpo_security" "gpo_sec" {
#     gpo_container = ad_gpo.gpo.id
#     password_policies {
#         minimum_password_length = 3
#     }
# 
#     system_services {
#         service_name = "TapiSrv"
#         startup_mode = "2"
#         acl = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;LA)"
#     }
#     
#     system_services {
#         service_name = "CertSvc"
#         startup_mode = "2"
#         acl = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;SY)(A;;CCLCSWLOCRRC;;;IU)S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)"
#     }
# 
# }
# 
