// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import (
	"testing"
)

func TestWriteEventLogPolicy(t *testing.T) {
	data := map[string]interface{}{
		"maximum_log_size": "10",
	}

	out, err := NewEventLogPolicy(data)
	if err != nil {
		t.Error(err)
	}

	if out.MaximumLogSize != "10" {
		t.Errorf("mismatch: MaximumLogSize. Expected 10 found %q", out.MaximumLogSize)
	}
}
