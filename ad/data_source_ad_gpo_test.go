package ad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceADGPO_basic(t *testing.T) {
	envVars := []string{"TF_VAR_ad_domain_name", "TF_VAR_ad_gpo_name"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGPOConfigBasic(),
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

func testAccDatasourceADGPOConfigBasic() string {
	return `

	variable "ad_domain_name" {}
	variable "ad_gpo_name" {}

	resource "ad_gpo" "gpo" {
		name        = var.ad_gpo_name
		domain      = var.ad_domain_name
	}

	data "ad_gpo" "g" {
	    name = ad_gpo.gpo.name
	}
	`
}
