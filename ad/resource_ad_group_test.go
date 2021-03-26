package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "security"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
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

func TestAccGroup_categories(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "security"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "system"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
		},
	})
}

func TestAccGroup_scopes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "domainlocal", "security", "some description"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "universal", "security", "some description"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
			{
				Config: testAccGroupConfigBasic("yourdomain.com", "test group", "testgroup", "global", "security", "some description"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupExists("ad_group.g", "yourdomain.com", "testgroup", true),
				),
			},
		},
	})
}

func testAccGroupConfigBasic(domain, name, sam, scope, gtype string, desc string) string {
	return fmt.Sprintf(`
	variable "name" { default = %q }
	variable "sam_account_name" { default = %q }
	variable "scope" { default = %q }
	variable "category" { default = %q }
	variable "container" { default = "CN=Users,dc=yourdomain,dc=com" }
	variable "description" { default = %q }

	resource "ad_group" "g" {
		name = var.name
		sam_account_name = var.sam_account_name
		scope = var.scope
		category = var.category
		container = var.container
		description = var.description
 	}
`, name, sam, scope, gtype, desc)
}

func testAccGroupExists(name, domain, groupSAM string, expected bool) resource.TestCheckFunc {
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
