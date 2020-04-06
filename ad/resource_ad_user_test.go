package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/ldaphelper"
)

func TestAccUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists("ad_user.a", "yourdomain.com", "testuser", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "yourdomain.com", "testuser", true),
				),
			},
		},
	})
}

func testAccUserConfigBasic(domain, username, password string) string {
	domainDN := getDomainFromDNSDomain(domain)
	principalName := fmt.Sprintf("%s@%s", username, domain)
	return fmt.Sprintf(`
	variable domain_dn { default = "%s" }
	variable principal_name { default = "%s" }
	variable password { default = "%s" }
	variable samaccountname { default = "%s" }

	resource "ad_user" "a" {
		domain_dn = var.domain_dn
		principal_name = var.principal_name
		sam_account_name = var.samaccountname
		initial_password = var.password
		display_name = "Terraform Test User"		
 	}
`, domainDN, principalName, password, username)
}

func testAccUserExists(name, domain, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("%s key not found on the server", name)
		}
		domainDN := getDomainFromDNSDomain(domain)
		ldapConn := testAccProvider.Meta().(ProviderConf).LDAPConn
		u, err := ldaphelper.GetUserFromLDAP(ldapConn, rs.Primary.ID, domainDN)
		if err != nil {
			if strings.Contains(err.Error(), "No entries found for filter") && !expected {
				return nil
			}
			return err
		}

		if u.SAMAccountName != username {
			return fmt.Errorf("username from LDAP does not match expected username, %s != %s", u.SAMAccountName, username)
		}
		return nil
	}
}
