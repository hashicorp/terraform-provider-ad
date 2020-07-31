package ad

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADOU_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADOUBasic("testOU", "dc=yourdomain,dc=com", "true", "test ou"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_ou.ods", "name",
						"ad_ou.o", "name",
					),
				),
			},
		},
	})
}

func testAccDataSourceADOUBasic(name, path, protected, description string) string {
	return fmt.Sprintf(`
	variable name { default = %q }
	variable path { default = %q }
	variable protected { default = %q }
	variable description { default = %q }

	resource "ad_ou" "o" {
		name = var.name
		path = var.path
		description = var.description
		protected = var.protected
	}

	data "ad_ou" "ods" {
		dn = ad_ou.o.dn
	}
`, name, path, protected, description)
}
