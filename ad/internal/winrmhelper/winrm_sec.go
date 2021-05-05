package winrmhelper

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/gposec"
	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
	"gopkg.in/ini.v1"
)

// GetSecIniFromResource buiilds the contents of the security settings ini file based on the data of the
// resource.
func GetSecIniFromResource(d *schema.ResourceData, schemaKeys map[string]*schema.Schema) (*ini.File, error) {
	loadOpts := ini.LoadOptions{
		AllowBooleanKeys:         true,
		KeyValueDelimiterOnWrite: "=",
		KeyValueDelimiters:       "=",
		IgnoreInlineComment:      true,
	}
	iniFile := ini.Empty(loadOpts)
	cfg := gposec.NewSecuritySettings()

	err := iniFile.ReflectFrom(cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.PopulateSecuritySettings(d, iniFile)
	if err != nil {
		return nil, err
	}

	return iniFile, nil

}

// GetSecIniContents returns a byte array with the contents of the INF file
// encoded in UTF-8 (since we get the ouput via stdout).
func GetSecIniContents(client *winrm.Client, gpo *GPO, execLocally, passCredentials bool, username, password string) ([]byte, error) {
	gptPath := fmt.Sprintf("%s\\Machine\\Microsoft\\Windows NT\\SecEdit\\GptTmpl.inf", gpo.basePath)
	log.Printf("[DEBUG] Getting security settings inf from %s", gptPath)

	cmd := fmt.Sprintf(`Get-Content "%s"`, gptPath)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving contents of %q: %s", gptPath, err)
	}
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("command to retrieve contents of %q failed, stderr: %s, stdout: %s", gptPath, result.StdErr, result.Stdout)
	}

	iniBytes := []byte(result.Stdout)
	return iniBytes, nil
}

// GetSecIniFromHost returns a struct representing the data retrieved from the host.
func GetSecIniFromHost(client *winrm.Client, gpo *GPO, execLocally, passCredentials bool, username, password string) (*gposec.SecuritySettings, error) {

	iniBytes, err := GetSecIniContents(client, gpo, execLocally, passCredentials, username, password)
	if err != nil {
		return nil, err
	}
	iniFile, err := gposec.ParseIniFile(iniBytes, false)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ini file, error: %s", err)
	}
	return iniFile, nil
}

// UploadSecIni uploads the security settings ini to the correct folder of a GPO and updates
// the GPO's gpt.ini by incrementing the computer version by 1.
func UploadSecIni(conn *winrm.Client, cpConn *winrmcp.Winrmcp, gpo *GPO, iniFile *ini.File, execLocally, passCredentials bool, username, password string) error {
	ini.LineBreak = "\r\n"
	buf := bytes.NewBuffer([]byte{})
	iniLocation := fmt.Sprintf("%s\\Machine\\Microsoft\\Windows NT\\SecEdit\\GptTmpl.inf", gpo.basePath)
	_, err := iniFile.WriteTo(buf)
	if err != nil {
		return fmt.Errorf("error while loading security INF file to buffer, error: %s ", err)
	}
	err = cpConn.Write(iniLocation, buf)
	if err != nil {
		return fmt.Errorf("error while writing ini file to %q: %s", iniLocation, err)
	}
	cVer := gpo.computerVersion + 1

	err = gpo.SetGPOVersions(conn, cpConn, gpo.userVersion, cVer, execLocally, passCredentials, username, password)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSecIni removes the ini file from the host and updates the GPO's  gpt.ini by incrementing the
// computer version by 1.
func RemoveSecIni(conn *winrm.Client, cpConn *winrmcp.Winrmcp, gpo *GPO, execLocally, passCredentials bool, username, password string) error {
	gptPath := fmt.Sprintf("%s\\Machine\\Microsoft\\Windows NT\\SecEdit\\GptTmpl.inf", gpo.basePath)
	log.Printf("[DEBUG] Getting security settings inf from %s", gptPath)

	cmd := fmt.Sprintf(`Remove-Item "%s"`, gptPath)
	result, err := RunWinRMCommand(conn, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return fmt.Errorf("error while retrieving contents of %q: %s", gptPath, err)
	}

	if result.ExitCode != 0 {
		if !strings.Contains(result.StdErr, "ItemNotFoundException") {
			return fmt.Errorf("error while removing %q: %s", gptPath, err)
		}
	}

	cVer := gpo.computerVersion + 1
	err = gpo.SetGPOVersions(conn, cpConn, gpo.userVersion, cVer, execLocally, passCredentials, username, password)
	if err != nil {
		return err
	}
	return nil
}
