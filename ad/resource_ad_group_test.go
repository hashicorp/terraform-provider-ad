package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGroup_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_domain_name",
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_group_scope",
		"TF_VAR_ad_group_category",
		"TF_VAR_ad_group_description",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_global"), os.Getenv("TF_VAR_ad_group_category_security")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
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

func TestAccResourceADGroup_categories(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_domain_name",
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_group_scope",
		"TF_VAR_ad_group_category",
		"TF_VAR_ad_group_description",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_global"), os.Getenv("TF_VAR_ad_group_category_security")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
				),
			},
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_global"), os.Getenv("TF_VAR_ad_group_category_distribution")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
				),
			},
		},
	})
}

func TestAccResourceADGroup_scopes(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_domain_name",
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_group_scope",
		"TF_VAR_ad_group_category",
		"TF_VAR_ad_group_description",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists("ad_group.g", os.Getenv(""), false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_domainlocal"), os.Getenv("TF_VAR_ad_group_category_security")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
				),
			},
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_universal"), os.Getenv("TF_VAR_ad_group_category_security")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
				),
			},
			{
				Config: testAccResourceADGroupConfigBasic(os.Getenv("TF_VAR_ad_group_scope_global"), os.Getenv("TF_VAR_ad_group_category_security")),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists("ad_group.g", os.Getenv("TF_VAR_ad_group_sam"), true),
				),
			},
		},
	})
}

func testAccResourceADGroupConfigBasic(scope, gtype string) string {
	return fmt.Sprintf(`
	variable "ad_group_name" {}
	variable "ad_group_sam"{}
	variable "scope" { default = %q }
	variable "category" { default = %q }
	variable "ad_group_container"{}
	variable "ad_group_description"{}

	resource "ad_group" "g" {
		name = var.ad_group_name
		sam_account_name = var.ad_group_sam
		scope = var.scope
		category = var.category
		container = var.ad_group_container
		description = var.ad_group_description
 	}
`, scope, gtype)
}

func testAccResourceADGroupExists(name, groupSAM string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("%s key not found on the server", name)
		}
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		u, err := winrmhelper.GetGroupFromHost(client, rs.Primary.ID, false)
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

		if u.SAMAccountName != groupSAM {
			return fmt.Errorf("username from LDAP does not match expected username, %s != %s", u.SAMAccountName, groupSAM)
		}

		if u.Scope != rs.Primary.Attributes["scope"] {
			return fmt.Errorf("actual scope does not match expected scope, %s != %s", rs.Primary.Attributes["scope"], u.Scope)
		}

		if u.Category != rs.Primary.Attributes["category"] {
			return fmt.Errorf("actual category does not match expected scope, %s != %s", rs.Primary.Attributes["category"], u.Category)
		}
		return nil
	}
}
