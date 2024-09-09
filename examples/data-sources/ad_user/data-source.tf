# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "ad_user" "u" {
    user_id = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
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

data "ad_user" "u2" {
    user_id = "CN=Test User,OU=Users,DC=contoso,DC=com"
}

output "testuser_guid" {
    value = data.ad_user.u2.id
}