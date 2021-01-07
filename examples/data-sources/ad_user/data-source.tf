data "ad_user" "u" {
    guid = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "username" {
    value = data.ad_user.u.sam_account_name
}

output "country" {
    value = data.ad_user.u.country
}

output "trusted_for_delegation" {
    value = data.ad_user.u.trusted_for_delegation
}