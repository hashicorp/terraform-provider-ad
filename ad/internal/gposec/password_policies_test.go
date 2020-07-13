package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWritePasswordPolicies(t *testing.T) {
	data := map[string]interface{}{
		"maximum_password_age": "10",
	}

	out := NewSecuritySettings()
	err := WritePasswordPolicies(data, out)
	if err != nil {
		t.Error(err)
	}

	if out.PasswordPolicies == nil {
		t.Errorf("PasswordPolicies struct is nil.")
		t.FailNow()
	}

	if out.PasswordPolicies.MaximumPasswordAge != "10" {
		t.Errorf("mismatch: MaximumPasswordAge. Expected 10 found %q", out.PasswordPolicies.MaximumPasswordAge)
	}
}

func TestPasswordPoliciesSetResourceData(t *testing.T) {

	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	pp := PasswordPolicies{
		MaximumPasswordAge: "10",
	}
	err := pp.SetResourceData("password_policies", d)
	if err != nil {
		t.Errorf("error while setting resource data: %s", err)
	}

	mpa := d.Get("password_policies.0.maximum_password_age").(string)

	if mpa != "10" {
		t.Errorf("unexpected value of password_policies.0.maximum_password_age. Expected 10 got %q", mpa)
	}

}
