package ad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceADComputer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADComputerBasic("testcomputer"),
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

func testAccDataSourceADComputerBasic(name string) string {
	return fmt.Sprintf(`
	variable "name" { default = %q }
	
	resource "ad_computer" "c" {
		name = var.name
	}	
	
	data "ad_computer" "dsc" {
		guid = ad_computer.c.guid
	}
	
	`, name)
}
