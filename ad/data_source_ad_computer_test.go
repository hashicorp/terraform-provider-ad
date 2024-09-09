// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADComputer_basic(t *testing.T) {
	envVars := []string{"TF_VAR_ad_computer_name"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADComputerBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_computer.dsc", "guid",
						"ad_computer.c", "guid",
					),
				),
			},
		},
	})
}

func testAccDataSourceADComputerBasic() string {
	return `

	variable "ad_computer_name" {}

	resource "ad_computer" "c" {
		name = var.ad_computer_name
	}	
	
	data "ad_computer" "dsc" {
		guid = ad_computer.c.guid
	}
	
	`
}
