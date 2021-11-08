package ad

import (
	"fmt"
	"strings"
	"testing"
	"os"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func TestAccResourceADGroupMember_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_group2_name",
		"TF_VAR_ad_group2_sam",
		"TF_VAR_ad_group2_container",
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupMemberExists("ad_group_member.gm", false, ""),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupMemberConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupMemberExists("ad_group_member.gm", true, os.Getenv("TF_VAR_ad_user_principal_name")),
				),
			},
			{
				ResourceName:      "ad_group_member.gm",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceADGroupMember_Update(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_group_name",
		"TF_VAR_ad_group_sam",
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_group2_name",
		"TF_VAR_ad_group2_sam",
		"TF_VAR_ad_group2_container",
		"TF_VAR_ad_group3_name",
		"TF_VAR_ad_group3_sam",
		"TF_VAR_ad_group3_container",
		"TF_VAR_ad_user_display_name",
		"TF_VAR_ad_user_sam",
		"TF_VAR_ad_user_password",
		"TF_VAR_ad_user_principal_name",
		"TF_VAR_ad_user_container",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupMemberExists("ad_group_member.gm", false, ""),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupMemberUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupMemberExists("ad_group_member.gm", true, os.Getenv("TF_VAR_ad_group2_name")),
				),
			},
		},
	})
}
func testAccResourceADGroupMemberExists(resourceName string, expected bool, member string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("%s resource not found", resourceName)
		}

		toks := strings.Split(rs.Primary.ID, "_")
		gm, err := winrmhelper.NewGroupMembershipFromHost(testAccProvider.Meta().(*config.ProviderConf), toks[0])
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

     	if expected && gm.GroupMembers[0].Name!=member {
			return fmt.Errorf("actual member (%s) does not match the expected member (%s)", gm.GroupMembers[0].Name, member)
		}

		return nil
	}
}

func testAccResourceADGroupMemberConfigBasic() string {
	return `

		variable "ad_group_name" {}
		variable "ad_group_sam" {}
		variable "ad_group_container" {}

		variable "ad_group2_name" {}
		variable "ad_group2_sam" {}
		variable "ad_group2_container" {}

		variable "ad_user_display_name" {}
		variable "ad_user_principal_name" {}
		variable "ad_user_sam" {}
		variable "ad_user_password" {}
		variable "ad_user_container" {}

		resource ad_group "g" {
			name             = var.ad_group_name
			sam_account_name = var.ad_group_sam
			container        = var.ad_group_container
		}

		resource ad_group "g2" {
			name             = var.ad_group2_name
			sam_account_name = var.ad_group2_sam
			container        = var.ad_group2_container
		}

		resource ad_user "u" {
			display_name     = var.ad_user_display_name
			principal_name   = var.ad_user_principal_name
			sam_account_name = var.ad_user_sam
			initial_password = var.ad_user_password
			container        = var.ad_user_container
		}

		resource ad_group_member "gm" {
			group_id = ad_group.g.id
			group_member  =  ad_user.u.id
		}
	`
}

func testAccResourceADGroupMemberUpdate() string {
	return `
		variable "ad_group_name" {}
		variable "ad_group_sam" {}
		variable "ad_group_container" {}

		variable "ad_group2_name" {}
		variable "ad_group2_sam" {}
		variable "ad_group2_container" {}

		variable "ad_group3_name" {}
		variable "ad_group3_sam" {}
		variable "ad_group3_container" {}

		variable "ad_user_display_name" {}
		variable "ad_user_principal_name" {}
		variable "ad_user_sam" {}
		variable "ad_user_password" {}
		variable "ad_user_container" {}

		resource ad_group "g" {
			name             = var.ad_group_name
			sam_account_name = var.ad_group_sam
			container        = var.ad_group_container
		}

		resource ad_group "g2" {
			name             = var.ad_group2_name
			sam_account_name = var.ad_group2_sam
			container        = var.ad_group2_container
		}

		resource ad_group "g3" {
			name             = var.ad_group3_name
			sam_account_name = var.ad_group3_sam
			container        = var.ad_group3_container
		}


		resource ad_user "u" {
			display_name     = var.ad_user_display_name
			principal_name   = var.ad_user_principal_name
			sam_account_name = var.ad_user_sam
			initial_password = var.ad_user_password
			container        = var.ad_user_container
		}

		resource ad_group_member "gm" {
			group_id = ad_group.g.id
			group_member  = ad_group.g2.id
		}
`
}
