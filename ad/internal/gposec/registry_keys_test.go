package gposec

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
	"gopkg.in/ini.v1"
)

func TestRegistryKeysSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rk := RegistryKeys{
		Keys: []string{
			`HKLM\Some\Key,2,D:ACL;STRINGGOESHERE`,
		},
	}

	err := rk.SetResourceData("registry_keys", d)
	if err != nil {
		t.Error(err)
	}

	if d.Get("registry_keys") == nil {
		t.Error("registry_keys set is nil")
		t.FailNow()
	}

	rkSet := d.Get("registry_keys").(*schema.Set)
	if rkSet.Len() == 0 {
		t.Error("empty RegistryKeys set")
		t.FailNow()
	}

	rkItem := rkSet.List()[0].(map[string]interface{})
	registryKey := rkItem["key_name"].(string)
	propMode := rkItem["propagation_mode"].(string)
	acl := rkItem["acl"].(string)

	if registryKey != `HKLM\Some\Key` || propMode != "2" || acl != "D:ACL;STRINGGOESHERE" {
		t.Errorf(`unexpected values found. Expected HKLM\Some\Key,2,D:ACL;STRINGGOESHERE got %s,%s,%s`, registryKey, propMode, acl)
	}
}

func newRKFromResource() (*RegistryKeys, error) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rData := []map[string]interface{}{
		{
			"key_name":         `HKLM\Some\Key`,
			"propagation_mode": "2",
			"acl":              "D:ACL;STRINGGOESHERE",
		},
	}
	err := d.Set("registry_keys", rData)
	if err != nil {
		return nil, err
	}

	rkSection, err := NewRegistryKeysFromResource(d.Get("registry_keys"))
	if err != nil {
		return nil, err
	}
	rk := rkSection.(*RegistryKeys)

	return rk, nil
}

func TestNewRegistryKeysFromResource(t *testing.T) {
	rk, err := newRKFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(rk.Keys) == 0 {
		t.Errorf("empty RegistryKeys struct found")
		t.FailNow()
	}

	if rk.Keys[0] != `"HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`RegistryKeys structure did not contain expected value "HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE". actual value was :%s`, rk.Keys[0])
	}

}

func TestSetIniData(t *testing.T) {
	rk, err := newRKFromResource()
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

	err = rk.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section := iniFile.Section("Registry Keys")

	if section.Body() != `"HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`RegistryKeys section body did not match expected value "HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE". actual value: %s`, section.Body())
	}

	iniFile = ini.Empty(loadOpts)
	rk.Keys = []string{}
	err = rk.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section = iniFile.Section("Registry Keys")

	if section.Body() != "" {
		t.Errorf(`RegistryKeys section body was not empty. actual value: %s`, section.Body())
	}
}

func TestLoadRegistryKeysFromIni(t *testing.T) {

	cfg := NewSecuritySettings()

	iniData := `
	[Registry Keys]
	"HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE"
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

	err = LoadRegistryKeysFromIni("Registry Keys", iniFile, cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(cfg.RegistryKeys.Keys) != 1 {
		t.Errorf("Invalid number of Keys found in RegistryKeys struct. Expected 1 got %d", len(cfg.RegistryKeys.Keys))
		t.FailNow()
	}

	if cfg.RegistryKeys.Keys[0] != `"HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`RegistryKeys Key did not match expected value "HKLM\Some\Key",2,"D:ACL;STRINGGOESHERE". actual value: %s`, cfg.RegistryKeys.Keys)
	}

	err = LoadRegistryKeysFromIni("Not Registry Keys", iniFile, cfg)
	if !strings.Contains(err.Error(), "error while parsing section") {
		t.Error(err)
	}
}
