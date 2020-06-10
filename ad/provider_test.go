package ad

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ad": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {

}

func getDomainFromDNSDomain(dnsDomain string) string {
	toks := strings.Split(dnsDomain, ".")
	for idx, tok := range toks {
		toks[idx] = fmt.Sprintf("dc=%s", tok)
	}
	domainDN := strings.Join(toks, ",")
	return domainDN
}
