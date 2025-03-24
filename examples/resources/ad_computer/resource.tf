# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "name" { default = "test" }
variable "pre2kname" { default = "TEST" }

resource "ad_computer" "c" {
  name      = var.name
  pre2kname = var.pre2kname
}
