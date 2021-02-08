package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccADGroupMembership_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccADGroupMembershipExists("ad_group_membership.gm", false, 0),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccADGroupMembershipConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccADGroupMembershipExists("ad_group_membership.gm", true, 2),
				),
			},
			{
				ResourceName:      "ad_group_membership.gm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccADGroupMembership_Uodate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccADGroupMembershipExists("ad_group_membership.gm", false, 0),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccADGroupMembershipUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccADGroupMembershipExists("ad_group_membership.gm", true, 3),
				),
			},
		},
	})
}
func testAccADGroupMembershipExists(resourceName string, expected bool, desiredMemberCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("%s resource not found", resourceName)
		}
		client, err := testAccProvider.Meta().(ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(ProviderConf).ReleaseWinRMClient(client)
		toks := strings.Split(rs.Primary.ID, "_")
		gm, err := winrmhelper.NewGroupMembershipFromHost(client, toks[0], false)
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

		if len(gm.GroupMembers) != desiredMemberCount {
			return fmt.Errorf("group member actual count (%d) does not match the expected number of members (%d)", len(gm.GroupMembers), desiredMemberCount)
		}
		return nil
	}
}

func testAccADGroupMembershipConfigBasic() string {
	return `
		resource ad_group "g" {
			name             = "testGroup"
			sam_account_name = "testgroup"
			container        = "CN=Users,dc=yourdomain,dc=com"
		}

		resource ad_group "g2" {
			name             = "memberGroup"
			sam_account_name = "membergroup"
			container        = "CN=Users,dc=yourdomain,dc=com"
		}

		resource ad_user "u" {
			display_name     = "test user"
			principal_name   = "testUser"
			sam_account_name = "testUser"
			initial_password = "SomethingRandom1234!!"
			container        = "CN=Users,DC=yourdomain,DC=com"
		}

		resource ad_group_membership "gm" {
			group_id = ad_group.g.id
			group_members  = [ ad_group.g2.id,ad_user.u.id]
		}
	`
}

func testAccADGroupMembershipUpdate() string {
	return `
		resource ad_group "g" {
			name             = "testGroup"
			sam_account_name = "testgroup"
			container        = "CN=Users,dc=yourdomain,dc=com"
		}

		resource ad_group "g2" {
			name             = "memberGroup"
			sam_account_name = "membergroup"
			container        = "CN=Users,dc=yourdomain,dc=com"
		}

		resource ad_group "g3" {
			name             = "memberGroup1"
			sam_account_name = "membergroup2"
			container        = "CN=Users,dc=yourdomain,dc=com"
		}

		resource ad_user "u" {
			display_name     = "test user"
			principal_name   = "testUser"
			sam_account_name = "testUser"
			initial_password = "SomethingRandom1234!!"
			container        = "CN=Users,DC=yourdomain,DC=com"
		}

		resource ad_group_membership "gm" {
			group_id = ad_group.g.id
			group_members  = [ ad_group.g2.id,ad_user.u.id,ad_group.g3.id]
		}
`
}
