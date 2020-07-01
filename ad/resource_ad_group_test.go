package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"
)

func TestAccGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "security"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
			{
				ResourceName:      "ad_group.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupConfigBasic(domain, name, sam, scope, gtype string) string {
	domainDN := getDomainFromDNSDomain(domain)
	return fmt.Sprintf(`
	variable "domain_dn" { default = %q }
	variable "display_name" { default = %q }
	variable "sam_account_name" { default = %q }
	variable "scope" { default = %q }
	variable "type" { default = %q }

	resource "ad_group" "g" {
		domain_dn = var.domain_dn
		display_name = var.display_name
		sam_account_name = var.sam_account_name
		scope = var.scope
		type = var.type
 	}
`, domainDN, name, sam, scope, gtype)
}

func testAccGroupExists(name, domain, groupSAM string, expected bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("%s key not found on the server", name)
		}
		ldapConn := testAccProvider.Meta().(ProviderConf).LDAPConn
		u, err := ldaphelper.GetGroupFromLDAP(ldapConn, rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "No entries found for filter") && !expected {
				return nil
			}
			return err
		}

		if u.SAMAccountName != groupSAM {
			return fmt.Errorf("username from LDAP does not match expected username, %s != %s", u.SAMAccountName, groupSAM)
		}
		return nil
	}
}
