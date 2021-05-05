package winrmhelper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
	"gopkg.in/ini.v1"
)

// GPO describes a Group Policy container, used to hold group policies
type GPO struct {
	Name            string `json:"DisplayName"`
	ID              string `json:"Id"`
	DN              string `json:"Path"`
	Domain          string `json:"DomainName"`
	Description     string `json:"Description"`
	NumericStatus   int    `json:"GpoStatus"`
	Status          string
	computerVersion uint16
	userVersion     uint16
	basePath        string
	gptIni          *ini.File
}

func getGPOCmdByName(name string) string {
	return fmt.Sprintf("Get-GPO -Name %s", name)
}

func getGPOCmdByGUID(guid string) string {
	return fmt.Sprintf("Get-GPO -Guid %s", guid)
}

// GPOStatusMap is used to translate the GPO status from a numeric format the json output returns
// to the string format we need to use when updating a GPO
var GPOStatusMap = map[int]string{
	0: "AllSettingsDisabled",
	1: "UserSettingsDisabled",
	2: "ComputerSettingsDisabled",
	3: "AllSettingsEnabled",
}

// unmarshallGPO unmarshalls the incoming byte array containing JSON
// into a GPO structure.
func unmarshallGPO(input []byte) (*GPO, error) {
	var gpo GPO
	err := json.Unmarshal(input, &gpo)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall json document with error %q, document was: %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling json response: %s", err)
	}
	status, ok := GPOStatusMap[gpo.NumericStatus]
	if !ok {
		return nil, fmt.Errorf("unknown GPO status %d", gpo.NumericStatus)
	}
	gpo.Status = status
	return &gpo, nil
}

// GetGPOFromHost returns a GPO structure populated by data from the DC server
func GetGPOFromHost(conn *winrm.Client, name, guid string, execLocally, passCredentials bool, username, password string) (*GPO, error) {
	start := time.Now().Unix()
	var cmd string
	if name != "" {
		cmd = getGPOCmdByName(name)
	} else if guid != "" {
		cmd = getGPOCmdByGUID(guid)
	}
	result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally, passCredentials, username, password)
	if err != nil {
		return nil, err
	}
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("command exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}
	gpo, err := unmarshallGPO([]byte(result.Stdout))
	if err != nil {
		return nil, err
	}

	basePath, err := gpo.loadGPOBasePath(conn, execLocally, passCredentials, username, password)
	if err != nil {
		return nil, err
	}
	gpo.basePath = basePath

	err = gpo.loadGPTIni(conn, execLocally, passCredentials, username, password)
	if err != nil {
		return nil, err
	}

	err = gpo.loadGPOVersions(conn, basePath)
	if err != nil {
		return nil, err
	}

	end := time.Now().Unix()
	log.Printf("[DEBUG] GPO from host took %d seconds", end-start)
	return gpo, nil
}

// GetGPOFromResource returns a GPO structure popuplated by data from TF
func GetGPOFromResource(d *schema.ResourceData) *GPO {
	g := GPO{
		Name:        SanitiseTFInput(d, "name"),
		Domain:      SanitiseTFInput(d, "domain"),
		Description: SanitiseTFInput(d, "description"),
		Status:      SanitiseTFInput(d, "status"),
		ID:          d.Id(),
	}
	return &g
}

// Rename renames a GPO to the given name
func (g *GPO) Rename(client *winrm.Client, target string, execLocally, passCredentials bool, username, password string) error {
	if g.ID == "" {
		return fmt.Errorf("gpo guid required")
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("Rename-GPO -Guid %s -TargetName %s", g.ID, g.Name))

	if g.Domain != "" {
		cmds = append(cmds, fmt.Sprintf("-Domain %s", g.Domain))
	}
	cmd := strings.Join(cmds, " ")
	_, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return err
	}
	return nil
}

//ChangeStatus Changes the status of a GPO
func (g *GPO) ChangeStatus(client *winrm.Client, status string, execLocally, passCredentials bool, username, password string) error {
	cmd := fmt.Sprintf(`(%s).GpoStatus = "%s"`, getGPOCmdByGUID(g.ID), status)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("status update failed with a non zero exit code (%d) stdout: %s stderr:%s",
			result.ExitCode, result.Stdout, result.StdErr)
	}

	return nil
}

// NewGPO uses Powershell over WinRM to create a script
func (g *GPO) NewGPO(client *winrm.Client, execLocally, passCredentials bool, username string, password string) (string, error) {

	if g.Name == "" {
		return "", fmt.Errorf("gpo name required")
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("New-GPO -Name %q", g.Name))

	if g.Domain != "" {
		cmds = append(cmds, fmt.Sprintf("-Domain %q", g.Domain))
	}

	if g.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Comment %q", g.Description))
	}

	result, err := RunWinRMCommand(client, cmds, true, false, execLocally, passCredentials, username, password)
	if err != nil {
		return "", err
	}
	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		if strings.Contains(result.StdErr, "GpoWithNameAlreadyExists") {
			return "", fmt.Errorf("there is another GPO named %q", g.Name)
		}
		return "", fmt.Errorf("command exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}
	gpo, err := unmarshallGPO([]byte(result.Stdout))
	if err != nil {
		return "", err
	}
	return gpo.ID, nil
}

// DeleteGPO delete the GPO container
func (g *GPO) DeleteGPO(client *winrm.Client, execLocally bool, passCredentials bool, username, password string) error {
	cmd := fmt.Sprintf("Remove-GPO -Name %s -Domain %s", g.Name, g.Domain)
	_, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		// Check if the resource is already deleted
		if strings.Contains(err.Error(), "GpoWithNameNotFound") {
			return nil
		}
		return err
	}
	return nil
}

// UpdateGPO updates the GPO container
func (g *GPO) UpdateGPO(client *winrm.Client, d *schema.ResourceData, execLocally, passCredentials bool, username, password string) (string, error) {
	if d.HasChange("name") {
		err := g.Rename(client, SanitiseTFInput(d, "name"), execLocally, passCredentials, username, password)
		if err != nil {
			return "", err
		}
	}

	if d.HasChange("status") {
		err := g.ChangeStatus(client, SanitiseTFInput(d, "status"), execLocally, passCredentials, username, password)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

// getGPOFilePath retrieves the AD Object of a GPO via powershell and returns the gPCFileSysPath
// property. This property points at the UNC that the GPO stores its configuration. We use the output
// of this function as well as GetsysVolPath to construct the GPO path on the DC's filesystem.
func (g *GPO) getGPOFilePath(client *winrm.Client, execLocally, passCredentials bool, username, password string) (string, error) {
	cmd := fmt.Sprintf("(Get-ADObject  -LDAPFilter '(&(objectClass=groupPolicyContainer)(cn={%s}))' -Properties gPCFilesysPath).gPCFilesysPath", g.ID)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return "", fmt.Errorf("error while retrieving GPO with %q path: %s", g.ID, err)
	}
	if result.ExitCode != 0 {
		return "", fmt.Errorf("error while retrieving SYSVOL path, stderr: %s, stdout: %s", result.StdErr, result.Stdout)
	}
	return result.Stdout, nil
}

//getSysVolPath returns the local path for the SYSVOL share on a Domain Controller. The combination of this
// and the value we get from getGPOFilePath is used to construct the GPO path on the DC's filesystem.
func getSysVolPath(client *winrm.Client, execLocally, passCredentials bool, username, password string) (string, error) {
	cmd := "(Get-SmbShare sysvol).path"
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return "", fmt.Errorf("error while retrieving SYSVOL path")
	}
	if result.ExitCode != 0 {
		return "", fmt.Errorf("error while retrieving SYSVOL path, stderr: %s, stdout: %s", result.StdErr, result.Stdout)
	}
	return result.Stdout, nil
}

// GetGPOBasePath returns the base path of a GPO on the DC. All GPO related files go
// in that location.
func (g *GPO) loadGPOBasePath(client *winrm.Client, execLocally, passCredentials bool, username, password string) (string, error) {
	gpoPath, err := g.getGPOFilePath(client, execLocally, passCredentials, username, password)
	if err != nil {
		return "", err
	}
	// gpoPath is a UNC. The first bit is the hostname and the second the share name
	// We are interested for the rest
	gPath := strings.Join(strings.Split(gpoPath, "\\")[4:], "\\")
	sysvolPath, err := getSysVolPath(client, execLocally, passCredentials, username, password)
	if err != nil {
		return "", err
	}
	gpoFinalPath := fmt.Sprintf("%s\\%s", sysvolPath, gPath)
	return gpoFinalPath, err
}

// GetGPOVersions returns the GPO versions for user and machine
func (g *GPO) loadGPOVersions(client *winrm.Client, gpoPath string) error {
	gpoVersionString, err := g.gptIni.Section("General").GetKey("Version")
	if err != nil {
		return fmt.Errorf("error while reading version for GPO: %q", g.ID)
	}
	gpoVersion, err := strconv.ParseInt(gpoVersionString.String(), 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert gpo version %s to uint32: %s", gpoVersionString, err)
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(gpoVersion))
	g.userVersion = binary.LittleEndian.Uint16(buf[:2])
	g.computerVersion = binary.LittleEndian.Uint16(buf[2:])
	return nil
}

// SetADGPOVersions updates AD with the given versions for a GPO
func (g *GPO) SetADGPOVersions(client *winrm.Client, gpoVersion uint32, execLocally, passCredentials bool, username, password string) error {
	cmd := fmt.Sprintf("$o=(Get-ADObject  -LDAPFilter '(&(objectClass=groupPolicyContainer)(cn={%s}))' -Properties *);$o.VersionNumber=%d;Set-AdObject -Instance $o", g.ID, gpoVersion)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return fmt.Errorf("error while setting new version in AD for GPO %q: %s", g.ID, err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("command to set the version of GPO %q in AD failed, stderr: %s, stdout: %s", g.ID, result.StdErr, result.Stdout)
	}
	return nil
}

// SetINIGPOVersions update gpt.ini with the new version
func (g *GPO) SetINIGPOVersions(client *winrm.Client, cpConn *winrmcp.Winrmcp, gpoVersion uint32) error {
	gpoVersionString, err := g.gptIni.Section("General").GetKey("Version")
	gpoVersionString.SetValue(strconv.Itoa(int(gpoVersion)))
	if err != nil {
		return fmt.Errorf("error while setting new GPT version to %d", gpoVersion)
	}

	buf := bytes.NewBuffer([]byte{})
	_, err = g.gptIni.WriteTo(buf)
	if err != nil {
		return fmt.Errorf("error while loading ini file contents in buffer")
	}

	gptPath := fmt.Sprintf("%s\\gpt.ini", g.basePath)
	err = cpConn.Write(gptPath, buf)
	if err != nil {
		return fmt.Errorf("error while writing ini file to %q: %s", gptPath, err)
	}

	return nil
}

// SetGPOVersions updates gpt.ini on the DC with the given values for user and computer version of a GPO.
func (g *GPO) SetGPOVersions(client *winrm.Client, cpConn *winrmcp.Winrmcp, userVersion, computerVersion uint16, execLocally, passCredentials bool, username, password string) error {
	outBuf := make([]byte, 4)
	binary.LittleEndian.PutUint16(outBuf[:2], computerVersion)
	binary.LittleEndian.PutUint16(outBuf[2:], userVersion)
	newVersion := binary.LittleEndian.Uint32(outBuf)

	err := g.SetINIGPOVersions(client, cpConn, newVersion)
	if err != nil {
		return err
	}

	err = g.SetADGPOVersions(client, newVersion, execLocally, passCredentials, username, password)
	if err != nil {
		return err
	}
	return nil
}

func (g *GPO) loadGPTIni(client *winrm.Client, execLocally, passCredentials bool, username, password string) error {
	gptPath := fmt.Sprintf("%s\\gpt.ini", g.basePath)
	log.Printf("[DEBUG] Getting GPT ini from %s", gptPath)
	cmd := fmt.Sprintf(`Get-Content "%s"`, gptPath)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally, passCredentials, username, password)
	if err != nil {
		return fmt.Errorf("error while retrieving contents of %q: %s", gptPath, err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("command to retrieve contents of %q failed, stderr: %s, stdout: %s", gptPath, result.StdErr, result.Stdout)
	}

	iniFile, err := ini.Load([]byte(result.Stdout))
	if err != nil {
		return fmt.Errorf("contents of %q are not an ini file: %s", gptPath, err)
	}

	// counting for "DEFAULT" and "General"
	if len(iniFile.Sections()) != 2 {
		return fmt.Errorf("found more than 1 sections in %q, aborting (Sections found: %#v)", gptPath, iniFile.SectionStrings())
	}

	// initialise version if not present.
	if !iniFile.Section("General").HasKey("Version") {
		_, err := iniFile.Section("General").NewKey("Version", "0")
		if err != nil {
			return fmt.Errorf("error while adding Version key to General section: %s", err)
		}
	}
	g.gptIni = iniFile

	return nil
}
