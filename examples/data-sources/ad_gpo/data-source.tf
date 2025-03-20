# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "ad_gpo" "g" {
    name = "Some GPO"
}

output "gpo_uuid" {
    value = data.ad_gpo.g.guid
}
