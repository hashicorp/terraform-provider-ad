package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGPO_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGPOExists("ad_gpo.gpo", "", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOConfigBasic("yourdomain.com", "tfgpo", "TF managed GPO", "AllSettingsEnabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", "tfgpo", true),
				),
			},
			{
				Config: testAccResourceADGPOConfigBasic("yourdomain.com", "tfgpo123", "TF managed GPO", "AllSettingsEnabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", "tfgpo123", true),
				),
			},
			{
				Config: testAccResourceADGPOConfigBasic("yourdomain.com", "tfgpo123", "TF managed GPO", "AllSettingsDisabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", "tfgpo123", true),
				),
			},
			{
				ResourceName:      "ad_gpo.gpo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceADGPOConfigBasic(domain, name, description, status string) string {
	return fmt.Sprintf(`

	variable "domain"      { default = "%s" }
	variable "name"        { default = "%s" }
	variable "description" { default = "%s" }
	variable "status"      { default = "%s" }

	resource "ad_gpo" "gpo" {
		name        = var.name
		domain      = var.domain
		description = var.description
		status      = var.status
	}
	`, domain, name, description, status)
}

func testAccResourceADGPOExists(resourceName, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}
		guid := rs.Primary.ID
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		gpo, err := winrmhelper.GetGPOFromHost(client, "", guid)
		if err != nil {
			// Check that the err is really because the GPO was not found
			// and not because of other issues
			if strings.Contains(err.Error(), "GpoWithIdNotFound") && !expected {
				return nil
			}
			return err
		}
		if name != gpo.Name {
			return fmt.Errorf("gpo name %q does not match expected name %q", gpo.Name, name)
		}
		return nil
	}
}
