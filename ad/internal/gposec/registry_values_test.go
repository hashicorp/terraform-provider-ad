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

func TestRegistryValuesSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rk := RegistryValues{
		Values: []string{
			`HKLM\Some\Key,2,keyvalue`,
		},
	}

	err := rk.SetResourceData("registry_values", d)
	if err != nil {
		t.Error(err)
	}

	if d.Get("registry_values") == nil {
		t.Error("registry_values set is nil")
		t.FailNow()
	}

	rkSet := d.Get("registry_values").(*schema.Set)
	if rkSet.Len() == 0 {
		t.Error("empty registry_values set")
		t.FailNow()
	}

	rkItem := rkSet.List()[0].(map[string]interface{})
	registryKey := rkItem["key_name"].(string)
	propMode := rkItem["value_type"].(string)
	value := rkItem["value"].(string)

	if registryKey != `HKLM\Some\Key` || propMode != "2" || value != "keyvalue" {
		t.Errorf(`unexpected values found. Expected HKLM\Some\Key,2,keyvalue got %s,%s,%s`, registryKey, propMode, value)
	}
}

func newRVFromResource() (*RegistryValues, error) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rData := []map[string]interface{}{
		{
			"key_name":   `HKLM\Some\Key`,
			"value_type": "2",
			"value":      "keyvalue",
		},
	}
	err := d.Set("registry_values", rData)
	if err != nil {
		return nil, err
	}

	rvSection, err := NewRegistryValuesFromResource(d.Get("registry_values"))
	if err != nil {
		return nil, err
	}
	rv := rvSection.(*RegistryValues)

	return rv, nil
}

func TestNewRegistryValuesFromResource(t *testing.T) {
	rv, err := newRVFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(rv.Values) == 0 {
		t.Errorf("empty RegistryValues struct found")
		t.FailNow()
	}

	if rv.Values[0] != `"HKLM\Some\Key",2,"keyvalue"` {
		t.Errorf(`RegistryValues structure did not contain expected value "HKLM\Some\Key",2,"keyvalue". actual value was :%s`, rv.Values[0])
	}

}

func TestRegistryValuesSetIniData(t *testing.T) {
	rv, err := newRVFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}
	iniFile := ini.Empty(loadOpts)
	ini.LineBreak = "\r\n"

	err = rv.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section := iniFile.Section("Registry Values")
	if section.Body() != `"HKLM\Some\Key",2,"keyvalue"` {
		t.Errorf(`RegistryValues section body did not match expected value "HKLM\Some\Key",2,"keyvalue". actual value: %s`, section.Body())
	}

	iniFile = ini.Empty(loadOpts)
	rv.Values = []string{}
	err = rv.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section = iniFile.Section("Registry Values")

	if section.Body() != "" {
		t.Errorf(`RegistryValues section body was not empty. actual value: %s`, section.Body())
	}
}

func TestLoadRegistryValuesFromIni(t *testing.T) {

	cfg := NewSecuritySettings()

	iniData := `
	[Registry Values]
	"HKLM\Some\Key",2,"keyvalue"
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

	err = LoadRegistryValuesFromIni("Registry Values", iniFile, cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(cfg.RegistryValues.Values) != 1 {
		t.Errorf("Invalid number of Values found in RegistryValues struct. Expected 1 got %d", len(cfg.RegistryValues.Values))
		t.FailNow()
	}

	if cfg.RegistryValues.Values[0] != `"HKLM\Some\Key",2,"keyvalue"` {
		t.Errorf(`RegistryValues Key did not match expected value "HKLM\Some\Key",2,"keyvalue". actual value: %s`, cfg.RegistryValues.Values)
	}

	err = LoadRegistryValuesFromIni("Not Registry Values", iniFile, cfg)
	if !strings.Contains(err.Error(), "error while parsing section") {
		t.Error(err)
	}
}
