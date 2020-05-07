package gposec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
	"gopkg.in/ini.v1"
)

func TestFileSystemSetResourceData(t *testing.T) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	fs := FileSystem{
		Paths: []string{
			"C:\\whatever,2,D:ACL;STRINGGOESHERE",
		},
	}

	fs.SetResourceData("filesystem", d)

	if d.Get("filesystem") == nil {
		t.Error("filesystem set is nil")
		t.FailNow()
	}

	fsSet := d.Get("filesystem").(*schema.Set)
	if fsSet.Len() == 0 {
		t.Error("empty filesystem set")
		t.FailNow()
	}

	fsItem := fsSet.List()[0].(map[string]interface{})
	path := fsItem["path"].(string)
	propMode := fsItem["propagation_mode"].(string)
	acl := fsItem["acl"].(string)

	if path != "C:\\whatever" || propMode != "2" || acl != "D:ACL;STRINGGOESHERE" {
		t.Errorf(`unexpected values found. Expected C:\\whatever,2,D:ACL;STRINGGOESHERE" got %s,%s,%s`, path, propMode, acl)
	}
}

func newFSFromResource() (*FileSystem, error) {
	r := schema.Resource{}
	r.Schema = adschema.GpoSecuritySchema()
	d := r.TestResourceData()

	rData := []map[string]interface{}{
		{
			"path":             `C:\whatever`,
			"propagation_mode": "2",
			"acl":              "D:ACL;STRINGGOESHERE",
		},
	}
	err := d.Set("filesystem", rData)
	if err != nil {
		return nil, err
	}

	fsSection, err := NewFileSystemFromResource(d.Get("filesystem"))
	if err != nil {
		return nil, err
	}
	fs := fsSection.(*FileSystem)

	return fs, nil
}

func TestNewFileSystemFromResource(t *testing.T) {
	fs, err := newFSFromResource()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(fs.Paths) == 0 {
		t.Errorf("empty filesystem struct found")
		t.FailNow()
	}

	if fs.Paths[0] != `"C:\whatever",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`Filesystem structure did not contain expected value "C:\whatever",2,"D:ACL;STRINGGOESHERE". actual value was :%s`, fs.Paths[0])
	}

}

func TestFileSystemSetIniData(t *testing.T) {
	fs, err := newFSFromResource()
	if err != nil {
		t.Log("WTF")
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

	err = fs.SetIniData(iniFile)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	section := iniFile.Section("File Security")

	if section.Body() != `"C:\whatever",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`filesystem section body did not match expected value "C:\whatever",2,"D:ACL;STRINGGOESHERE". actual value: %s`, section.Body())
	}
}

func TestLoadFileSystemFromIni(t *testing.T) {
	cfg := NewSecuritySettings()

	iniData := `
	[File Security]
	"C:\whatever",2,"D:ACL;STRINGGOESHERE"
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

	err = LoadFileSystemFromIni("File Security", iniFile, cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(cfg.FileSystem.Paths) != 1 {
		t.Errorf("Invalid number of Paths found in FileSystem struct. Expected 1 got %d", len(cfg.FileSystem.Paths))
		t.FailNow()
	}

	if cfg.FileSystem.Paths[0] != `"C:\whatever",2,"D:ACL;STRINGGOESHERE"` {
		t.Errorf(`Filesystem Path did not match expected value "C:\whatever",2,"D:ACL;STRINGGOESHERE". actual value: %s`, cfg.FileSystem.Paths)
	}
}
