package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADOU_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_ou_name",
		"TF_VAR_ad_ou_description",
		"TF_VAR_ad_ou_path",
		"TF_VAR_ad_ou_protected",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADOUExists("ad_ou.o", "", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADOUConfigBasic("", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", os.Getenv("TF_VAR_ad_ou_name"), true),
				),
			},
			{
				Config: testAccResourceADOUConfigBasic("-renamed", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", fmt.Sprintf("%s-renamed",
						os.Getenv("TF_VAR_ad_ou_name")), true),
				),
			},
			{
				Config: testAccResourceADOUConfigBasic("", false),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", os.Getenv("TF_VAR_ad_ou_name"), true),
				),
			},
			{
				ResourceName:      "ad_ou.o",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceADOUConfigBasic(nameSuffix string, protection bool) string {
	return fmt.Sprintf(`
variable ad_ou_name {}
variable ad_ou_path {}
variable ad_ou_description {}
variable protected { default = %t}

resource "ad_ou" "o" { 
    name = "${var.ad_ou_name}%s"
    path = var.ad_ou_path
    description = var.ad_ou_description
    protected = var.protected
}
`, protection, nameSuffix)
}

func testAccResourceADOUExists(resource, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}
		guid := rs.Primary.ID
		ou, err := winrmhelper.NewOrgUnitFromHost(testAccProvider.Meta().(*config.ProviderConf), guid, "", "")
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}
		if ou.Name != name {
			return fmt.Errorf("OU name %q does not match expected name %q", ou.Name, name)
		}
		return nil

	}
}
