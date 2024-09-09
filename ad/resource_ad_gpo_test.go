// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func TestAccResourceADGPO_basic(t *testing.T) {
	envVars := []string{
		"TF_VAR_ad_gpo_name",
		"TF_VAR_ad_gpo_domain",
		"TF_VAR_ad_gpo_description",
		"TF_VAR_ad_gpo_status",
	}

	gpoName := os.Getenv("TF_VAR_ad_gpo_name")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGPOExists("ad_gpo.gpo", gpoName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", gpoName, true),
				),
			},
			{
				Config: testAccResourceADGPOConfigBasic("-renamed"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", fmt.Sprintf("%s-renamed", gpoName), true),
				),
			},
			{
				Config: testAccResourceADGPOConfigBasic(""),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists("ad_gpo.gpo", gpoName, true),
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

func testAccResourceADGPOConfigBasic(suffix string) string {
	return fmt.Sprintf(`

	variable "ad_gpo_domain" {}
	variable "ad_gpo_name" {}
	variable "ad_gpo_description" {}
	variable "ad_gpo_status" {}

	resource "ad_gpo" "gpo" {
		name        = "${var.ad_gpo_name}%s"
		domain      = var.ad_gpo_domain
		description = var.ad_gpo_description
		status      = var.ad_gpo_status
	}
	`, suffix)
}

func testAccResourceADGPOExists(resourceName, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}
		guid := rs.Primary.ID
		client, err := testAccProvider.Meta().(*config.ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(*config.ProviderConf).ReleaseWinRMClient(client)

		gpo, err := winrmhelper.GetGPOFromHost(testAccProvider.Meta().(*config.ProviderConf), "", guid)
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
