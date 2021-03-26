data "ad_gmsa" "g" {
    gmsa_id = "testGmsa$"
}

output "sam" {
    value = data.ad_gmsa.g.sam_account_name
}

output "trusted_for_delegation" {
    value = data.ad_gmsa.g.trusted_for_delegation
}

data "ad_gmsa" "g2" {
    gmsa_id = "CN=testGmsa,OU=gmsas,DC=contoso,DC=com"
}

output "testgmsa_guid" {
    value = data.ad_gmsa.g2.id
}