# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "domain" { default = "yourdomain.com" }
variable "name" { default = "tfGPO" }

resource "ad_gpo" "gpo" {
  name   = var.name
  domain = var.domain
}
