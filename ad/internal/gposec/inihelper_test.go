package gposec

import (
	"testing"
)

func TestStartIniFile(t *testing.T) {
	iniFile := NewSecuritySettings()
	if iniFile.Unicode.Unicode != "yes" || iniFile.Version.Revision != 1 || iniFile.Version.Signature != "\"$CHICAGO$\"" {
		t.Errorf("inifile does not contain the basic headers or values: %#v", iniFile)
	}

}
