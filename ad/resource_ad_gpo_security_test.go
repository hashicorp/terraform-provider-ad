package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGPOSecurity_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_gpo_name",
		"TF_VAR_ad_gpo_domain",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, envVars) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccResourceADGPOSecurityExists("ad_gpo_security.gpo_sec", false)),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOSecurityConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOSecurityExists("ad_gpo_security.gpo_sec", true),
				),
			},
			{
				ResourceName:      "ad_gpo_security.gpo_sec",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceADGPOSecurityExists(resourceName string, desired bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}

		toks := strings.Split(rs.Primary.ID, "_")
		if len(toks) != 2 {
			return fmt.Errorf("resource ID %q does not match <guid>_securitysettings", rs.Primary.ID)
		}
		guid := toks[0]
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		gpo, err := winrmhelper.GetGPOFromHost(client, "", guid, false)
		if err != nil {
			// if the GPO got destroyed first then the rest of the entities depending on it
			// are also destroyed.
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}

		_, err = winrmhelper.GetSecIniFromHost(client, gpo, false)
		if err != nil {
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccResourceADGPOSecurityConfigBasic() string {
	return `
variable "ad_gpo_domain" {}
variable "ad_gpo_name" {}

resource "ad_gpo" "gpo" {
    name        = var.ad_gpo_name
    domain      = var.ad_gpo_domain
}

resource "ad_gpo_security" "gpo_sec" {
    gpo_container = ad_gpo.gpo.id
   password_policies {
        minimum_password_length = 3
    }
}
`
}
