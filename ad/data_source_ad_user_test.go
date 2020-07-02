package ad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceADUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADUserBasic("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo"),
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

func testAccDataSourceADUserBasic(domain, username, password string) string {
	principalName := fmt.Sprintf("%s@%s", username, domain)
	return fmt.Sprintf(`
	variable "principal_name" { default = %q }
	variable "password" { default = %q }
	variable "samaccountname" { default = %q }

	resource "ad_user" "a" {
		principal_name = var.principal_name
		sam_account_name = var.samaccountname
		initial_password = var.password
		display_name = "Terraform Test User"
		container = "CN=Users,DC=yourdomain,DC=com"	
	}
	 
	 data "ad_user" "d" {
		 user_dn = ad_user.a.id
	 }
`, principalName, password, username)
}
