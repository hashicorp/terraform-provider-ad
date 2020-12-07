package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADComputer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists("ad_computer.c", "testcomputer", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigBasic("testcomputer", "TESTCOMPUTER$"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", "testcomputer", true),
				),
			},
			{
				ResourceName:      "ad_computer.c",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceADComputer_move(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists("ad_computer.c", "testcomputer", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigBasic("testcomputer", "TESTCOMPUTER$"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", "testcomputer", true),
				),
			},
			{
				Config: testAccResourceADComputerConfigMove("testcomputer", "TESTCOMPUTER$"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", "testcomputer", true),
				),
			},
		},
	})
}

func testAccResourceADComputerConfigBasic(name, prewin2kname string) string {
	return fmt.Sprintf(`
variable "name" { default = %q }
variable "pre2kname" { default = %q }

resource "ad_computer" "c" {
	name = var.name
	pre2kname = var.pre2kname
}
`, name, prewin2kname)
}

func testAccResourceADComputerConfigMove(name, prewin2kname string) string {
	return fmt.Sprintf(`
variable "name" { default = %q }
variable "pre2kname" { default = %q }

resource "ad_ou" "o" { 
	name = "anotherou"
	path = "dc=yourdomain,dc=com"
}
resource "ad_computer" "c" {
	name = var.name
	pre2kname = var.pre2kname
	container = ad_ou.o.dn
}
`, name, prewin2kname)
}

func testAccResourceADComputerExists(resource, name string, expected bool) resource.TestCheckFunc {
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
		computer, err := winrmhelper.NewComputerFromHost(client, guid)
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}

		if computer.Name != name {
			return fmt.Errorf("Computer name %q does not match expected name %q", computer.Name, name)
		}
		return nil
	}
}
