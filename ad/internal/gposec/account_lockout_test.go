package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteAccountLockout(t *testing.T) {
	data := map[string]interface{}{
		"force_logoff_when_hour_expire": "10",
	}

	out := NewSecuritySettings()
	err := WriteAccountLockout(data, out)
	if err != nil {
		t.Error(err)

	}
	if out.AccountLockout == nil {
		t.Errorf("AccountLockout struct is nil.")
		t.FailNow()
	}

	if out.AccountLockout.ForceLogoffWhenHourExpire != "10" {
		t.Errorf("unexpected value of ForceLogOffWhenHourExpire. Expected 10 found %q", out.AccountLockout.ForceLogoffWhenHourExpire)
	}
}

func TestALSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	pp := AccountLockout{
		LockoutBadCount: "10",
	}
	pp.SetResourceData("account_lockout", d)

	mpa := d.Get("account_lockout.0.lockout_bad_count").(string)

	if mpa != "10" {
		t.Errorf("unexpected value of account_lockout.0.lockout_bad_count_age. Expected 10 got %q", mpa)
	}
}
