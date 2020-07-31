package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", true),
				),
			},
			{
				ResourceName:            "ad_user.a",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password"},
			},
		},
	})
}

func TestAccUser_modify(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser123", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", true),
				),
			},
			{
				Config: testAccUserConfigBasic("dc=yourdomain,dc=com", "testuser123", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser123", true),
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
			testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigUAC("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo", "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "dc=yourdomain,dc=com", false, false),
				),
			},
			{
				Config: testAccUserConfigUAC("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo", "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "dc=yourdomain,dc=com", true, false),
				),
			},
			{
				Config: testAccUserConfigUAC("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo", "false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "dc=yourdomain,dc=com", false, true),
				),
			},
			{
				Config: testAccUserConfigUAC("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo", "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC("ad_user.a", "dc=yourdomain,dc=com", true, true),
				),
			},
		},
	})
}

func defaultVariablesSection(domain, username, password string) string {
	principalName := fmt.Sprintf("%s@%s", username, domain)
	return fmt.Sprintf(`
	variable "principal_name" { default = %q }
	variable "password" { default = %q }
	variable "samaccountname" { default = %q }

	`, principalName, password, username)
}

func defaultUserSection() string {
	return `
	principal_name = var.principal_name
	sam_account_name = var.samaccountname
	initial_password = var.password
	display_name = "Terraform Test User"
	container = "CN=Users,DC=yourdomain,DC=com"	
	`
}
func testAccUserConfigBasic(domain, username, password string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s    		
 	}`, defaultVariablesSection(domain, username, password), defaultUserSection())

}

func testAccUserConfigUAC(domain, username, password, enabled, expires string) string {
	return fmt.Sprintf(`%s
	variable "enabled" { default = %q }
	variable "password_never_expires" { default = %q }

	resource "ad_user" "a" {%s
		enabled = var.enabled
		password_never_expires = var.password_never_expires
 	}
`, defaultVariablesSection(domain, username, password), enabled, expires, defaultUserSection())
}

func retrieveADUserFromRunningState(name, domain string, s *terraform.State) (*winrmhelper.User, error) {
	rs, ok := s.RootModule().Resources[name]

	if !ok {
		return nil, fmt.Errorf("%s key not found in state", name)
	}
	client := testAccProvider.Meta().(ProviderConf).WinRMClient
	u, err := winrmhelper.GetUserFromHost(client, rs.Primary.ID)

	return u, err

}

func testAccUserExists(name, domain, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, domain, s)
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
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

func testCheckADUserUAC(name, domain string, enabledState, passwordNeverExpires bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, domain, s)

		if err != nil {
			return err
		}

		if u.Enabled != enabledState {
			return fmt.Errorf("enabled state in AD did not match expected value. AD: %t, expected: %t", u.Enabled, enabledState)
		}

		if u.PasswordNeverExpires != passwordNeverExpires {
			return fmt.Errorf("password_never_expires state in AD did not match expected value. AD: %t, expected: %t", u.PasswordNeverExpires, enabledState)
		}
		return nil
	}
}
