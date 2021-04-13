package ad

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ad": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T, envVars []string) {
	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val == "" {
			t.Fatalf("%s must be set for acceptance tests to work", envVar)
		}
	}
}
