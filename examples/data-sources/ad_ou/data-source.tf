# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "ad_ou" "g" {
    dn = "OU=SomeOU,dc=yourdomain,dc=com"
}

output "ou_uuid" {
    value = data.ad_ou.g.id
}
