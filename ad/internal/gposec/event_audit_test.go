package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteEventAudit(t *testing.T) {
	data := map[string]interface{}{
		"audit_logon_events": "1",
	}

	out := NewSecuritySettings()
	err := WriteEventAudit(data, out)
	if err != nil {
		t.Error(err)
	}

	if out.EventAudit == nil {
		t.Errorf("EventAudit struct is nil.")
		t.FailNow()
	}

	if out.EventAudit.AuditLogonEvents != "1" {
		t.Errorf("mismatch: MaximumLogSize. Expected 10 found %q", out.AuditLog.MaximumLogSize)
	}
}

func TestEventAuditSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	al := EventAudit{
		AuditAccountLogon: "10",
	}
	al.SetResourceData("event_audit", d)

	mls := d.Get("event_audit.0.audit_account_logon").(string)

	if mls != "10" {
		t.Errorf("unexpected value of event_audit.0.audit_account_logon. Expected 10 got %q", mls)
	}
}
