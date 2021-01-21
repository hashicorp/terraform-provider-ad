package ad

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

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
			{
				Config: testAccUserConfigAttributes("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", true),
				),
			},
			{
				Config: testAccUserConfigBasic("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", true),
				),
			},
		},
	})
}

func TestAccUser_custom_attributes_basic(t *testing.T) {
	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	username := "testuser"
	password := "thu2too'W?ieJ}a^g0zo"
	domainDN := "dc=yourdomain,dc=com"
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists(resourceName, domainDN, username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigCustomAttributes(domainDN, username, password, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, domainDN, caConfig),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password", "custom_attributes"},
			},
		},
	})
}

func TestAccUser_custom_attributes_extended(t *testing.T) {
	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	caConfig2 := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"], "comment": "another string"}`
	username := "testuser"
	password := "thu2too'W?ieJ}a^g0zo"
	domainDN := "dc=yourdomain,dc=com"
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccUserExists(resourceName, domainDN, username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigBasic("dc=yourdomain,dc=com", "testuser", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserExists("ad_user.a", "dc=yourdomain,dc=com", "testuser", true),
				),
			},
			{
				Config: testAccUserConfigCustomAttributes(domainDN, username, password, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, domainDN, caConfig),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password", "custom_attributes"},
			},
			{
				Config: testAccUserConfigCustomAttributes(domainDN, username, password, caConfig2),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, domainDN, caConfig2),
				),
			},
			{
				Config: testAccUserConfigCustomAttributes(domainDN, username, password, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, domainDN, caConfig),
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
			{
				Config: testAccUserConfigMoved("dc=yourdomain,dc=com", "testuser123", "thu2too'W?ieJ}a^g0zo"),
				Check: resource.ComposeTestCheckFunc(
					testAccUserContainer("ad_user.a", "dc=yourdomain,dc=com", "ou=newOU,DC=yourdomain,DC=com"),
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

func defaultUserSection(container string) string {
	if container == "" {
		container = `"CN=Users,DC=yourdomain,DC=com"`
	}
	return fmt.Sprintf(`
	principal_name = var.principal_name
	sam_account_name = var.samaccountname
	initial_password = var.password
	display_name = "Terraform Test User"
	container = %s
	`, container)
}

func testAccUserConfigBasic(domain, username, password string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s
 	}`, defaultVariablesSection(domain, username, password), defaultUserSection(""))

}

func testAccUserConfigAttributes(domain, username, password string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s
	  city                      = "City"
	  company                   = "Company"
	  country                   = "us"
	  department                = "Department"
	  description               = "Description"
	  division                  = "Division"
	  email_address             = "some@email.com"
	  employee_id               = "id"
	  employee_number           = "number"
	  fax                       = "Fax"
	  given_name                = "GivenName"
	  home_directory            = "HomeDirectory"
	  home_drive                = "HomeDrive"
	  home_phone                = "HomePhone"
	  home_page                 = "HomePage"
	  initials                  = "Initia"
	  mobile_phone              = "MobilePhone"
	  office                    = "Office"
	  office_phone              = "OfficePhone"
	  organization              = "Organization"
	  other_name                = "OtherName"
	  po_box                    = "POBox"
	  postal_code               = "PostalCode"
	  state                     = "State"
	  street_address            = "StreetAddress"
	  surname                   = "Surname"
	  title                     = "Title"
	  smart_card_logon_required = false
	  trusted_for_delegation    = true
	}`, defaultVariablesSection(domain, username, password), defaultUserSection(""))

}

func testAccUserConfigCustomAttributes(domain, username, password, customAttributes string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s
		custom_attributes = jsonencode(%s)
 	}`,
		defaultVariablesSection(domain, username, password),
		defaultUserSection(""),
		customAttributes)
}

func testAccUserConfigMoved(domain, username, password string) string {
	return fmt.Sprintf(`%s

	resource "ad_ou" "o" {
		name = "newOU"
		path = "DC=yourdomain,DC=com"
		description = "ou for user move test"
		protected = false
	}

	resource "ad_user" "a" {%s
 	}`, defaultVariablesSection(domain, username, password), defaultUserSection("ad_ou.o.dn"))

}

func testAccUserConfigUAC(domain, username, password, enabled, expires string) string {
	return fmt.Sprintf(`%s
	variable "enabled" { default = %q }
	variable "password_never_expires" { default = %q }

	resource "ad_user" "a" {%s
		enabled = var.enabled
		password_never_expires = var.password_never_expires
 	}
`, defaultVariablesSection(domain, username, password), enabled, expires, defaultUserSection(""))
}

func retrieveADUserFromRunningState(name, domain string, s *terraform.State, attributeList []string) (*winrmhelper.User, error) {
	rs, ok := s.RootModule().Resources[name]

	if !ok {
		return nil, fmt.Errorf("%s key not found in state", name)
	}
	client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return nil, err
	}
	defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, rs.Primary.ID, attributeList)

	return u, err

}

func testAccUserContainer(name, domain, expectedContainer string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		u, err := retrieveADUserFromRunningState(name, domain, s, nil)
		if err != nil {
			return err
		}

		if strings.EqualFold(u.Container, expectedContainer) {
			return fmt.Errorf("user container mismatch: expected %q found %q", u.Container, expectedContainer)
		}
		return nil
	}
}

func testAccUserExists(name, domain, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, domain, s, nil)
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
		u, err := retrieveADUserFromRunningState(name, domain, s, nil)

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

func testCheckADUserCustomAttribute(name, domain, customAttributes string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ca, err := structure.ExpandJsonFromString(customAttributes)
		if err != nil {
			return err
		}

		attributeList := []string{}
		for k := range ca {
			attributeList = append(attributeList, k)
		}

		u, err := retrieveADUserFromRunningState(name, domain, s, attributeList)
		if err != nil {
			return err
		}

		sortedCA := winrmhelper.SortInnerSlice(ca)
		sortedStateCA := winrmhelper.SortInnerSlice(u.CustomAttributes)

		if !reflect.DeepEqual(sortedCA, sortedStateCA) {
			return fmt.Errorf("attributes %#v returned from host do not match the attributes defined in the configuration: %#v vs %#v", attributeList, ca, u.CustomAttributes)
		}
		return nil
	}
}
