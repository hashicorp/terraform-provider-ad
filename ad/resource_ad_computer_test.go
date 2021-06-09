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

func TestAccResourceADComputer_basic(t *testing.T) {
	computerName := os.Getenv("TF_VAR_ad_computer_name")

	envVars := []string{"TF_VAR_ad_computer_name", "TF_VAR_ad_computer_sam"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists("ad_computer.c", computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", computerName, true),
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

func TestAccResourceADComputer_description(t *testing.T) {
	computerName := os.Getenv("TF_VAR_ad_computer_name")
	description := os.Getenv("TF_VAR_ad_computer_description")

	envVars := []string{"TF_VAR_ad_computer_name", "TF_VAR_ad_computer_description", "TF_VAR_ad_computer_sam"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerDescriptionExists("ad_computer.c", computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", computerName, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigDescription(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerDescriptionExists("ad_computer.c", description, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerDescriptionExists("ad_computer.c", "", true),
				),
			},
		},
	})
}

func TestAccResourceADComputer_move(t *testing.T) {
	computerName := os.Getenv("TF_VAR_ad_computer_name")

	envVars := []string{"TF_VAR_ad_computer_name", "TF_VAR_ad_computer_sam", "TF_VAR_ad_ou_name", "TF_VAR_ad_ou_path"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists("ad_computer.c", computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", computerName, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigMove(),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists("ad_computer.c", computerName, true),
				),
			},
		},
	})
}

func testAccResourceADComputerConfigBasic() string {
	return `
variable "ad_computer_name" {}
variable "ad_computer_sam" {}

resource "ad_computer" "c" {
	name = var.ad_computer_name
	pre2kname = var.ad_computer_sam
}
`
}

func testAccResourceADComputerConfigDescription() string {
	return `
variable "ad_computer_name" {}
variable "ad_computer_sam" {}
variable "ad_computer_description" {}

resource "ad_computer" "c" {
	name = var.ad_computer_name
	pre2kname = var.ad_computer_sam
	description = var.ad_computer_description
}
`
}

func testAccResourceADComputerConfigMove() string {
	return `
variable "ad_computer_name" {}
variable "ad_computer_sam" {}
variable "ad_ou_name" {}
variable "ad_ou_path" {}

resource "ad_ou" "o" { 
	name = var.ad_ou_name
	path = var.ad_ou_path
}

resource "ad_computer" "c" {
	name = var.ad_computer_name
	pre2kname = var.ad_computer_sam
	container = ad_ou.o.dn
}
`
}

func testAccResourceADComputerExists(resource, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}

		guid := rs.Primary.ID
		computer, err := winrmhelper.NewComputerFromHost(testAccProvider.Meta().(*config.ProviderConf), guid)
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}

		if computer.Name != name {
			return fmt.Errorf("computer name %q does not match expected name %q", computer.Name, name)
		}
		return nil
	}
}

func testAccResourceADComputerDescriptionExists(resource, description string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}

		guid := rs.Primary.ID
		computer, err := winrmhelper.NewComputerFromHost(testAccProvider.Meta().(*config.ProviderConf), guid)
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}

		if computer.Description != description {
			return fmt.Errorf("computer description %q does not match expected description %q", computer.Description, description)
		}
		return nil
	}
}
