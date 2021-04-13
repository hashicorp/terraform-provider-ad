package ad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADUser_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_container",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_user.d", "id",
						"ad_user.a", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceADUserBasic() string {
	return `
	variable "ad_user_principal_name" {}
	variable "ad_user_password" {}
	variable "ad_user_sam" {}
	variable "ad_user_display_name" {}
	variable "ad_user_container" {}

	resource "ad_user" "a" {
		principal_name = var.ad_user_principal_name
		sam_account_name = var.ad_user_sam
		initial_password = var.ad_user_password
		display_name = var.ad_user_display_name
		container = var.ad_user_container
	}
	 
	 data "ad_user" "d" {
		 user_id = ad_user.a.id
	 }
`
}
