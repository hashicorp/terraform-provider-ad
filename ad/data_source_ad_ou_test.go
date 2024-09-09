// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADOU_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_ou_name",
		"TF_VAR_ad_ou_path",
		"TF_VAR_ad_ou_protected",
		"TF_VAR_ad_ou_description",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADOUBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_ou.ods", "name",
						"ad_ou.o", "name",
					),
				),
			},
		},
	})
}

func testAccDataSourceADOUBasic() string {
	return `
	variable ad_ou_name {}
	variable ad_ou_path {}
	variable ad_ou_protected {}
	variable ad_ou_description {}

	resource "ad_ou" "o" {
		name = var.ad_ou_name
		path = var.ad_ou_path
		description = var.ad_ou_description
		protected = var.ad_ou_protected
	}

	data "ad_ou" "ods" {
		dn = ad_ou.o.dn
	}
`
}
