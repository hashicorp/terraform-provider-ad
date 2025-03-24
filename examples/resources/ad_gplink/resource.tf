# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "ad_ou" "o" {
  name        = "gplinktestOU"
  path        = "dc=yourdomain,dc=com"
  description = "OU for gplink tests"
  protected   = false
}

resource "ad_gpo" "g" {
  name        = "gplinktestGPO"
  domain      = "yourdomain.com"
  description = "gpo for gplink tests"
  status      = "AllSettingsEnabled"
}

resource "ad_gplink" "og" {
  gpo_guid  = ad_gpo.g.id
  target_dn = ad_ou.o.dn
  enforced  = true
  enabled   = true
}
