package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGPOSecurity_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccResourceADGPOSecurityExists("ad_gpo_security.gpo_sec", false)),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOSecurityConfigBasic("yourdomain.com", "tfgpo"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOSecurityExists("ad_gpo_security.gpo_sec", true),
				),
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
		client := testAccProvider.Meta().(ProviderConf).WinRMClient
		gpo, err := winrmhelper.GetGPOFromHost(client, "", guid)
		if err != nil {
			// if the GPO got destroyed first then the rest of the entities depending on it
			// are also destroyed.
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}

		_, err = winrmhelper.GetSecIniFromHost(client, gpo)
		if err != nil {
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccResourceADGPOSecurityConfigBasic(domain, name string) string {
	return fmt.Sprintf(`
variable "domain"      { default = "%s" }
variable "name"        { default = "%s" }

resource "ad_gpo" "gpo" {
    name        = var.name
    domain      = var.domain
}

resource "ad_gpo_security" "gpo_sec" {
    gpo_container = ad_gpo.gpo.id
   password_policies {
        minimum_password_length = 3
    }
}
`, domain, name)
}
