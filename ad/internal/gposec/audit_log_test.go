// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
)

func TestWriteAuditLog(t *testing.T) {
	data := map[string]interface{}{
		"maximum_log_size": "10",
	}

	out := NewSecuritySettings()
	err := WriteAuditLog(data, out)
	if err != nil {
		t.Error(err)
	}

	if out.AuditLog == nil {
		t.Errorf("AuditLog struct is nil.")
		t.FailNow()
	}

	if out.AuditLog.MaximumLogSize != "10" {
		t.Errorf("mismatch: MaximumLogSize. Expected 10 found %q", out.AuditLog.MaximumLogSize)
	}
}

func TestAuditLogSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	al := AuditLog{
		EventLogPolicy: EventLogPolicy{MaximumLogSize: "10"},
	}
	err := al.SetResourceData("audit_log", d)
	if err != nil {
		t.Errorf("error while setting resource data: %s", err)
	}
	mls := d.Get("audit_log.0.maximum_log_size").(string)

	if mls != "10" {
		t.Errorf("unexpected value of audit_log.0.maximum_log_size. Expected 10 got %q", mls)
	}
}
