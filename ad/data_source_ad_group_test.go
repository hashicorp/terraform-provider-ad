package ad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceADGroup_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_scope",
		"TF_VAR_ad_group_category",
		"TF_VAR_ad_group_container",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGroupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_group.d", "id",
						"ad_group.g", "id",
					),
				),
			},
		},
	})

}

func testAccDatasourceADGroupConfigBasic() string {
	return `
	variable "ad_group_name" {}
	variable "ad_group_sam" {}
	variable "ad_group_scope" {}
	variable "ad_group_category" {}
	variable "ad_group_container" {}

	resource "ad_group" "g" {
		name = var.ad_group_name
		sam_account_name = var.ad_group_sam
		scope = var.ad_group_scope
		category = var.ad_group_category
		container = var.ad_group_container
	 }
	 
	 data "ad_group" "d" {
		 group_id = ad_group.g.id
	 }
`
}
