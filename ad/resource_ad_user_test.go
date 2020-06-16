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

func TestAccUser_modify(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists("ad_user.a", "yourdomain.com", "testuser123", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "yourdomain.com", "testuser", true),
				),
			},
			{
				Config: testAccUserConfigBasic("yourdomain.com", "testuser123", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "yourdomain.com", "testuser123", true),
				),
			},
		},
	})
}

func TestAccUser_UAC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists("ad_user.a", "yourdomain.com", "testuser", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigUAC("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo", "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "yourdomain.com", false, false),
				),
			},
			{
				Config: testAccUserConfigUAC("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo", "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "yourdomain.com", true, false),
				),
			},
			{
				Config: testAccUserConfigUAC("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo", "false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "yourdomain.com", false, true),
				),
			},
			{
				Config: testAccUserConfigUAC("yourdomain.com", "testuser", "thu2too'W?ieJ}a^g0zo", "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "yourdomain.com", true, true),
				),
			},
		},
	})
}

func defaultVariablesSection(domain, username, password string) string {
	domainDN := getDomainFromDNSDomain(domain)
	principalName := fmt.Sprintf("%s@%s", username, domain)
	return fmt.Sprintf(`
	variable "domain_dn" { default = "%s" }
	variable "principal_name" { default = "%s" }
	variable "password" { default = "%s" }
	variable "samaccountname" { default = "%s" }

	`, domainDN, principalName, password, username)

}

func defaultUserSection() string {
	return `
	domain_dn = var.domain_dn
	principal_name = var.principal_name
	sam_account_name = var.samaccountname
	initial_password = var.password
	display_name = "Terraform Test User"	
	`
}
func testAccUserConfigBasic(domain, username, password string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s    		
 	}`, defaultVariablesSection(domain, username, password), defaultUserSection())

}

func testAccUserConfigUAC(domain, username, password, disabled, expires string) string {
	return fmt.Sprintf(`%s
	variable "disabled" { default = "%s" }
	variable "password_never_expires" { default = "%s" }

	resource "ad_user" "a" {%s
		disabled = var.disabled
		password_never_expires = var.password_never_expires
 	}
`, defaultVariablesSection(domain, username, password), disabled, expires, defaultUserSection())
}

func retrieveADUserFromRunningState(name, domain string, s *terraform.State) (*ldaphelper.User, error) {
	rs, ok := s.RootModule().Resources[name]

	if !ok {
		return nil, fmt.Errorf("%s key not found in stater", name)
	}
	domainDN := getDomainFromDNSDomain(domain)
	ldapConn := testAccProvider.Meta().(ProviderConf).LDAPConn
	u, err := ldaphelper.GetUserFromLDAP(ldapConn, rs.Primary.ID, domainDN)

	return u, err

}

func testAccUserExists(name, domain, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, domain, s)
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

func testCheckADUserUAC(name, domain string, disabledState, passwordNeverExpires bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, domain, s)

		if err != nil {
			return err
		}

		if u.Disabled != disabledState {
			return fmt.Errorf("disabled state in AD did not match expected value. AD: %t, expected: %t", u.Disabled, disabledState)
		}

		if u.PasswordNeverExpires != passwordNeverExpires {
			return fmt.Errorf("password_never_expires state in AD did not match expected value. AD: %t, expected: %t", u.PasswordNeverExpires, disabledState)
		}
		return nil
	}
}
