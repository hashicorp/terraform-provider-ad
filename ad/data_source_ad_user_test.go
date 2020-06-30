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
	domainDN := getDomainFromDNSDomain(domain)
	principalName := fmt.Sprintf("%s@%s", username, domain)
	return fmt.Sprintf(`
	variable "domain_dn" { default = %q }
	variable "principal_name" { default = %q }
	variable "password" { default = %q }
	variable "samaccountname" { default = %q }

	resource "ad_user" "a" {
		domain_dn = var.domain_dn
		principal_name = var.principal_name
		sam_account_name = var.samaccountname
		initial_password = var.password
		display_name = "Terraform Test User"		
	 }
	 
	 data "ad_user" "d" {
		 domain_dn = var.domain_dn
		 user_dn = ad_user.a.id
	 }
`, domainDN, principalName, password, username)
}
