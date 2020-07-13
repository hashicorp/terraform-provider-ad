package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteApplicationLog(t *testing.T) {
	data := map[string]interface{}{
		"maximum_log_size": "10",
	}

	out := NewSecuritySettings()
	err := WriteApplicationLog(data, out)
	if err != nil {
		t.Error(err)
	}

	if out.ApplicationLog == nil {
		t.Errorf("ApplicationLog struct is nil.")
		t.FailNow()
	}

	if out.ApplicationLog.MaximumLogSize != "10" {
		t.Errorf("mismatch: MaximumLogSize. Expected 10 found %q", out.ApplicationLog.MaximumLogSize)
	}
}

func TestApplicationLogSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	al := ApplicationLog{
		EventLogPolicy: EventLogPolicy{MaximumLogSize: "10"},
	}
	err := al.SetResourceData("application_log", d)
	if err != nil {
		t.Errorf("error while setting resource data: %s", err)
	}
	mls := d.Get("application_log.0.maximum_log_size").(string)

	if mls != "10" {
		t.Errorf("unexpected value of application_log.0.maximum_log_size. Expected 10 got %q", mls)
	}
}
