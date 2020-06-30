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
	domainDN := getDomainFromDNSDomain(domain)
	return fmt.Sprintf(`
	variable "domain_dn" { default = "%s" }
	variable "display_name" { default = "%s" }
	variable "sam_account_name" { default = "%s" }
	variable "scope" { default = "%s" }
	variable "type" { default = "%s" }

	resource "ad_group" "g" {
		domain_dn = var.domain_dn
		display_name = var.display_name
		sam_account_name = var.sam_account_name
		scope = var.scope
		type = var.type
	 }
	 
	 data "ad_group" "d" {
		 domain_dn = var.domain_dn
		 dn = ad_group.g.id
	 }
`, domainDN, name, sam, scope, gtype)
}
