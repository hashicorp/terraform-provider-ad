// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gposec

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
	"gopkg.in/ini.v1"
)

func TestSystemServicesSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	svcs := SystemServices{
		Services: []string{
			`somesvc,2,D:ACL;STRINGGOESHERE`,
		},
	}

	err := svcs.SetResourceData("system_services", d)
	if err != nil {
		t.Error(err)
	}

	if d.Get("system_services") == nil {
		t.Error("system_services set is nil")
		t.FailNow()
	}

	svcsSet := d.Get("system_services").(*schema.Set)
	if svcsSet.Len() == 0 {
		t.Error("empty SystemServices set")
		t.FailNow()
	}

	svcItem := svcsSet.List()[0].(map[string]interface{})
	svcName := svcItem["service_name"].(string)
	svcMode := svcItem["startup_mode"].(string)
	svcACL := svcItem["acl"].(string)

	if svcName != "somesvc" || svcMode != "2" || svcACL != "D:ACL;STRINGGOESHERE" {
		t.Errorf("unexpected values found. Expected somesvc,2,D:ACL;STRINGGOESHERE got %s,%s,%s", svcName, svcMode, svcACL)
	}
}

func newSvcFromResource() (*SystemServices, error) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rData := []map[string]interface{}{
		{
			"service_name": `somesvc`,
			"startup_mode": "2",
			"acl":          "D:ACL;STRINGGOESHERE",
		},
	}
	err := d.Set("system_services", rData)
	if err != nil {
		return nil, err
	}

	svcsSection, err := NewSystemServicesFromResource(d.Get("system_services"))
	if err != nil {
		return nil, err
	}
	sSvc := svcsSection.(*SystemServices)

	return sSvc, nil
}

func TestNewSystemServicesFromResource(t *testing.T) {
	svcs, err := newSvcFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(svcs.Services) == 0 {
		t.Errorf("empty SystemServices struct found")
		t.FailNow()
	}

	if svcs.Services[0] != `"somesvc",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`SystemServices structure did not contain expected value "somesvc",2,"D:ACL;STRINGGOESHERE". actual value was: %s`, svcs.Services[0])
	}
}

func TestSystemServicesValuesFromIni(t *testing.T) {
	cfg := NewSecuritySettings()

	iniData := `
	[System Services]
	"somesvc",2,"D:ACL;STRINGGOESHERE"
	`

	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}
	iniFile, err := ini.LoadSources(loadOpts, []byte(iniData))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = LoadSystemServicesFromIni("System Services", iniFile, cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(cfg.SystemServices.Services) != 1 {
		t.Errorf("Invalid number of Services found in SystemServices struct. Expected 1 got %d", len(cfg.SystemServices.Services))
	}

	if cfg.SystemServices.Services[0] != `somesvc,2,D:ACL;STRINGGOESHERE` {
		t.Errorf(`SystemServices key did not match expected value "somesvc",2,"D:ACL;STRINGGOESHERE". actual value: %s`, cfg.SystemServices.Services[0])
	}

	err = LoadSystemServicesFromIni("Not System Services", iniFile, cfg)
	if !strings.Contains(err.Error(), "error while parsing section") {
		t.Error(err)
	}
}
