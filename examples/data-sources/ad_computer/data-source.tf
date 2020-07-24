data "ad_computer" "c" {
    dn = "cn=test123,cn=Computers,dc=yourdomain,dc=com"
}

output "computer_guid" {
  value = data.ad_computer.c.guid
}
