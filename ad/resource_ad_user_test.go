package ad

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADUser_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}

	username := os.Getenv("TF_VAR_ad_user_sam")
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, username, true),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password"},
			},
			{
				Config: testAccResourceADUserConfigAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, username, true),
				),
			},
			{
				Config: testAccResourceADUserConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, username, true),
				),
			},
		},
	})
}

func TestAccResourceADUser_custom_attributes_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}

	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	username := os.Getenv("TF_VAR_ad_user_sam")
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigCustomAttributes(caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
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

func TestAccResourceADUser_custom_attributes_extended(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}

	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	caConfig2 := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"], "comment": "another string"}`
	username := os.Getenv("TF_VAR_ad_user_sam")
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, username, true),
				),
			},
			{
				Config: testAccResourceADUserConfigCustomAttributes(caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password", "custom_attributes"},
			},
			{
				Config: testAccResourceADUserConfigCustomAttributes(caConfig2),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig2),
				),
			},
			{
				Config: testAccResourceADUserConfigCustomAttributes(caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
				),
			},
		},
	})
}

func TestAccResourceADUser_modify(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_ou_name",
		"TF_VAR_ad_ou_description",
		"TF_VAR_ad_ou_path",
		"TF_VAR_ad_ou_protected",
	}

	username := os.Getenv("TF_VAR_ad_user_sam")
	usernameSuffix := "renamed"
	renamedUsername := fmt.Sprintf("%s%s", username, usernameSuffix)
	resourceName := "ad_user.a"
	expectedContainerDN := fmt.Sprintf("%s,%s", os.Getenv("TF_VAR_ad_ou_name"),
		os.Getenv("TF_VAR_ad_ou_path"))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, renamedUsername, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, username, true),
				),
			},
			{
				Config: testAccResourceADUserConfigBasic(usernameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, renamedUsername, true),
				),
			},
			{
				Config: testAccResourceADUserConfigMoved(usernameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserContainer(resourceName, expectedContainerDN),
				),
			},
		},
	})
}

func TestAccResourceADUser_UAC(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}
	username := os.Getenv("TF_VAR_ad_user_sam")
	resourceName := "ad_user.a"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists("ad_user.a", username, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigUAC("false", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, false, false),
				),
			},
			{
				Config: testAccResourceADUserConfigUAC("true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, true, false),
				),
			},
			{
				Config: testAccResourceADUserConfigUAC("false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, false, true),
				),
			},
			{
				Config: testAccResourceADUserConfigUAC("true", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, true, true),
				),
			},
		},
	})
}

func defaultVariablesSection() string {
	return `
	variable "ad_user_principal_name"  {}
	variable "ad_user_password" {}
	variable "ad_user_sam" {}
	variable "ad_user_display_name" {}
	`
}

func defaultUserSection(usernameSuffix, container string) string {
	return fmt.Sprintf(`
	principal_name = var.ad_user_principal_name
	sam_account_name = "${var.ad_user_sam}%s"
	initial_password = var.ad_user_password
	display_name = var.ad_user_display_name
	container = %s
	`, usernameSuffix, container)
}

func testAccResourceADUserConfigBasic(usernameSuffix string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s
 	}`, defaultVariablesSection(), defaultUserSection(usernameSuffix, fmt.Sprintf("%q",
		os.Getenv("TF_VAR_ad_user_container"))))

}

func testAccResourceADUserConfigAttributes() string {
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
	}`, defaultVariablesSection(), defaultUserSection("", fmt.Sprintf("%q",
		os.Getenv("TF_VAR_ad_user_container"))))

}

func testAccResourceADUserConfigCustomAttributes(customAttributes string) string {
	return fmt.Sprintf(`%s
	resource "ad_user" "a" {%s
		custom_attributes = jsonencode(%s)
 	}`,
		defaultVariablesSection(),
		defaultUserSection("", fmt.Sprintf("%q", os.Getenv("TF_VAR_ad_user_container"))),
		customAttributes)
}

func testAccResourceADUserConfigMoved(usernameSuffix string) string {
	return fmt.Sprintf(`%s
	variable ad_ou_name {}
	variable ad_ou_path {}
	variable ad_ou_description {}
	variable ad_ou_protected {}
	
	resource "ad_ou" "o" { 
		name = var.ad_ou_name
		path = var.ad_ou_path
		description = var.ad_ou_description
		protected = var.ad_ou_protected
	}

	resource "ad_user" "a" {%s
 	}`, defaultVariablesSection(), defaultUserSection(usernameSuffix, "ad_ou.o.dn"))

}

func testAccResourceADUserConfigUAC(enabled, expires string) string {
	return fmt.Sprintf(`%s
	variable "enabled" { default = %q }
	variable "password_never_expires" { default = %q }

	resource "ad_user" "a" {%s
		enabled = var.enabled
		password_never_expires = var.password_never_expires
 	}
`, defaultVariablesSection(), enabled, expires, defaultUserSection("",
		fmt.Sprintf("%q", os.Getenv("TF_VAR_ad_user_container"))))
}

func retrieveADUserFromRunningState(name string, s *terraform.State, attributeList []string) (*winrmhelper.User, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("%s key not found in state", name)
	}
	u, err := winrmhelper.GetUserFromHost(testAccProvider.Meta().(*config.ProviderConf), rs.Primary.ID, attributeList)

	return u, err

}

func testAccResourceADUserContainer(name, expectedContainer string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		u, err := retrieveADUserFromRunningState(name, s, nil)
		if err != nil {
			return err
		}

		if strings.EqualFold(u.Container, expectedContainer) {
			return fmt.Errorf("user container mismatch: expected %q found %q", u.Container, expectedContainer)
		}
		return nil
	}
}

func testAccResourceADUserExists(name, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, s, nil)
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

func testCheckADUserUAC(name string, enabledState, passwordNeverExpires bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, s, nil)

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

func testCheckADUserCustomAttribute(name, customAttributes string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ca, err := structure.ExpandJsonFromString(customAttributes)
		if err != nil {
			return err
		}

		attributeList := []string{}
		for k := range ca {
			attributeList = append(attributeList, k)
		}

		u, err := retrieveADUserFromRunningState(name, s, attributeList)
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
