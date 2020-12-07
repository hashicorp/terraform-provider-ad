package ad

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGPLink_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
	//lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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

	resource "ad_ou" "o" { 
		name = "gplinktestOU"
		path = "dc=yourdomain,dc=com"
		description = "OU for gplink tests"
		protected = false
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

	resource "ad_ou" "o" { 
		name = "gplinktestOU"
		path = "dc=yourdomain,dc=com"
		description = "OU for gplink tests"
		protected = false
	}
		
	resource "ad_gpo" "g" {
		name        = "gplinktestGPO"
		domain      = "yourdomain.com"
		description = "gpo for gplink tests"
		status      = "AllSettingsEnabled"
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
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		idParts := strings.SplitN(id, "_", 2)
		if len(idParts) != 2 {
			return fmt.Errorf("malformed ID for GPLink resource with ID %q", id)
		}
		gplink, err := winrmhelper.GetGPLinkFromHost(client, idParts[0], idParts[1])
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
