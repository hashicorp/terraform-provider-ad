package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteKerberosSettings(t *testing.T) {
	data := map[string]interface{}{
		"max_service_age": "10",
	}

	out := NewSecuritySettings()
	err := WriteKerberosPolicy(data, out)
	if err != nil {
		t.Error(err)

	}
	if out.KerberosPolicy == nil {
		t.Errorf("KerberosSettings struct is nil.")
		t.FailNow()
	}

	if out.KerberosPolicy.MaxServiceAge != "10" {
		t.Errorf("mismatch: MaxServiceAge. Expected 10 found %q", out.KerberosPolicy.MaxServiceAge)
	}
}

func TestKerberosPolicySetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	al := KerberosPolicy{
		MaxTicketAge: "10",
	}
	al.SetResourceData("kerberos_policy", d)

	mls := d.Get("kerberos_policy.0.max_ticket_age").(string)

	if mls != "10" {
		t.Errorf("unexpected value of kerberos_policy.0.max_ticket_age. Expected 10 got %q", mls)
	}
}
