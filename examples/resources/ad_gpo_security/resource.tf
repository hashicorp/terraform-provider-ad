# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable domain { default = "yourdomain.com" }
variable gpo_name { default = "TestGPO" }

resource "ad_gpo" "gpo" {
  name   = var.gpo_name
  domain = var.domain
}

resource "ad_gpo_security" "gpo_sec" {
  gpo_container = ad_gpo.gpo.id
  password_policies {
    minimum_password_length = 3
  }

  system_services {
    service_name = "TapiSrv"
    startup_mode = "2"
    acl          = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;LA)"
  }

  system_services {
    service_name = "CertSvc"
    startup_mode = "2"
    acl          = "D:AR(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;SY)(A;;CCLCSWLOCRRC;;;IU)S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)"
  }

}

