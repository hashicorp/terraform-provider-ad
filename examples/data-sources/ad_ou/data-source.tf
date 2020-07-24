data "ad_ou" "g" {
    dn = "OU=SomeOU,dc=yourdomain,dc=com"
}

output "ou_uuid" {
    value = data.ad_ou.g.id
}
