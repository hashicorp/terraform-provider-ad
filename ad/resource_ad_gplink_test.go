package ad

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGPLink_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_ou_name",
		"TF_VAR_ad_ou_path",
		"TF_VAR_ad_ou_protected",
		"TF_VAR_ad_ou_description",
		"TF_VAR_ad_gpo_name",
		"TF_VAR_ad_gpo_domain",
		"TF_VAR_ad_gpo_description",
		"TF_VAR_ad_gpo_status",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGPLinkExists("ad_gplink.og", 1, true, true, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPLinkConfigBasic(true, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists("ad_gplink.og", 1, true, true, true),
				),
			},
			{
				ResourceName:      "ad_gplink.og",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigBasic(true, false, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists("ad_gplink.og", 1, true, false, true),
				),
			},
			{
				ResourceName:      "ad_gplink.og",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigBasic(false, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists("ad_gplink.og", 1, false, true, true),
				),
			},
			{
				ResourceName:      "ad_gplink.og",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigBasic(false, false, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists("ad_gplink.og", 1, false, false, true),
				),
			},
			{
				ResourceName:      "ad_gplink.og",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceADGPLink_badguid(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_ou_name",
		"TF_VAR_ad_ou_path",
		"TF_VAR_ad_ou_protected",
		"TF_VAR_ad_ou_description",
	}
	//lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceADGPLinkConfigBadGUID(false, false, 1),
				ExpectError: regexp.MustCompile("is not a valid uuid"),
			},
		},
	})
}

func testAccResourceADGPLinkConfigBadGUID(enforced, enabled bool, order int) string {
	return fmt.Sprintf(`
	variable ad_ou_name {}
	variable ad_ou_path {}
	variable ad_ou_protected {}
	variable ad_ou_description {}

	resource "ad_ou" "o" {
		name = var.ad_ou_name
		path = var.ad_ou_path
		description = var.ad_ou_description
		protected = var.ad_ou_protected
	}
		
	resource "ad_gplink" "og" { 
		gpo_guid = "something-horribly-wrong"
		target_dn = ad_ou.o.dn
		enforced = %t
		enabled = %t
		order = %d
	}
	
	`, enforced, enabled, order)
}

func testAccResourceADGPLinkConfigBasic(enforced, enabled bool, order int) string {
	return fmt.Sprintf(`
	variable ad_ou_name {}
	variable ad_ou_path {}
	variable ad_ou_protected {}
	variable ad_ou_description {}
	variable ad_gpo_name {}
	variable ad_gpo_domain {}
	variable ad_gpo_description {}
	variable ad_gpo_status {}

	resource "ad_ou" "o" {
		name = var.ad_ou_name
		path = var.ad_ou_path
		description = var.ad_ou_description
		protected = var.ad_ou_protected
	}
		
	resource "ad_gpo" "g" {
		name        = var.ad_gpo_name
		domain      = var.ad_gpo_domain
		description = var.ad_gpo_description
		status      = var.ad_gpo_status
	}
	
	resource "ad_gplink" "og" { 
		gpo_guid = ad_gpo.g.id
		target_dn = ad_ou.o.dn
		enforced = %t
		enabled = %t
		order = %d
	}
	
	`, enforced, enabled, order)
}

func testAccResourceADGPLinkExists(resourceName string, order int, enforced, enabled, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}
		id := rs.Primary.ID

		idParts := strings.SplitN(id, "_", 2)
		if len(idParts) != 2 {
			return fmt.Errorf("malformed ID for GPLink resource with ID %q", id)
		}
		gplink, err := winrmhelper.GetGPLinkFromHost(testAccProvider.Meta().(*config.ProviderConf), idParts[0], idParts[1])
		if err != nil {
			// Check that the err is really because the GPO was not found
			// and not because of other issues
			if strings.Contains(err.Error(), "did not find") && !expected {
				return nil
			}
			return err
		}

		if gplink.Enabled != enabled {
			return fmt.Errorf("gplink enabled status (%t) does not match expected status (%t)", gplink.Enabled, enabled)
		}

		if gplink.Enforced != enforced {
			return fmt.Errorf("gplink enforced status (%t) does not match expected status (%t)", gplink.Enforced, enforced)
		}

		if gplink.Order != order {
			return fmt.Errorf("gplink order %d does not match expected order %d", gplink.Order, order)
		}

		return nil
	}
}
