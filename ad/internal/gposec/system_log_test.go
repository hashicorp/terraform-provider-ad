package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteSystemLog(t *testing.T) {
	data := map[string]interface{}{
		"maximum_log_size": "10",
	}

	out := NewSecuritySettings()
	err := WriteSystemLog(data, out)
	if err != nil {
		t.Error(err)
	}

	if out.SystemLog == nil {
		t.Errorf("SystemLog struct is nil.")
		t.FailNow()
	}

	if out.SystemLog.MaximumLogSize != "10" {
		t.Errorf("mismatch: MaximumLogSize. Expected 10 found %q", out.SystemLog.MaximumLogSize)
	}
}

func TestSystemLogSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	al := SystemLog{
		EventLogPolicy: EventLogPolicy{MaximumLogSize: "10"},
	}
	err := al.SetResourceData("system_log", d)
	if err != nil {
		t.Errorf("error while setting resource data: %s", err)
	}

	mls := d.Get("system_log.0.maximum_log_size").(string)

	if mls != "10" {
		t.Errorf("unexpected value of system_log.0.maximum_log_size. Expected 10 got %q", mls)
	}
}
