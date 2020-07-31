package ad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceADGPO_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGPOConfigBasic("yourdomain.com", "tfgpo"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_gpo.g", "id",
						"ad_gpo.gpo", "id",
					),
				),
			},
		},
	})
}

func testAccDatasourceADGPOConfigBasic(domain, name string) string {
	return fmt.Sprintf(`

	variable"domain"      { default = "%s" }
	variable "name"        { default = "%s" }

	resource "ad_gpo" "gpo" {
		name        = var.name
		domain      = var.domain
	}

	data "ad_gpo" "g" {
	    name = ad_gpo.gpo.name
	}
	`, domain, name)
}
