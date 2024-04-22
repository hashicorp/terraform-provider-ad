# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "ad_ou" "o" { 
    name = "gplinktestOU"
    path = "dc=yourdomain,dc=com"
    description = "OU for gplink tests"
    protected = false
}
