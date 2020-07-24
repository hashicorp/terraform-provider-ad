package ad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDatasourceADGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "security"),
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

func testAccDatasourceADGroupConfigBasic(domain, name, sam, scope, gtype string) string {
	return fmt.Sprintf(`
	variable "name" { default = "%s" }
	variable "sam_account_name" { default = "%s" }
	variable "scope" { default = "%s" }
	variable "category" { default = "%s" }
	variable "container" { default = "cn=Users,dc=yourdomain,dc=com" }

	resource "ad_group" "g" {
		name = var.name
		sam_account_name = var.sam_account_name
		scope = var.scope
		category = var.category
		container = var.container
	 }
	 
	 data "ad_group" "d" {
		 guid = ad_group.g.id
	 }
`, name, sam, scope, gtype)
}
