package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADOU_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADOUExists("ad_ou.o", "", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADOUConfigBasic("testOU", "dc=yourdomain,dc=com", "some description", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", "testOU", true),
				),
			},
			{
				Config: testAccResourceADOUConfigBasic("testOU1", "dc=yourdomain,dc=com", "some description", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", "testOU1", true),
				),
			},
			{
				Config: testAccResourceADOUConfigBasic("testOU", "dc=yourdomain,dc=com", "some description", false),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists("ad_ou.o", "testOU", false),
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

func testAccResourceADOUConfigBasic(name, path, description string, protected bool) string {
	return fmt.Sprintf(`
variable name { default = "%s" }
variable path { default = "%s" }
variable description { default = "%s" }
variable protected { default = %t }

resource "ad_ou" "o" { 
    name = var.name
    path = var.path
    description = var.description
    protected = var.protected
}
`, name, path, description, protected)
}

func testAccResourceADOUExists(resource, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		guid := rs.Primary.ID
		ou, err := winrmhelper.NewOrgUnitFromHost(client, guid, "", "")
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
