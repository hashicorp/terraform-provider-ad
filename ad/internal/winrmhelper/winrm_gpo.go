package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/masterzen/winrm"
)

// GPO describes a Group Policy container, used to hold group policies
type GPO struct {
	Name          string `json:"DisplayName"`
	ID            string `json:"Id"`
	DN            string `json:"Path"`
	Domain        string `json:"DomainName"`
	Description   string `json:"Description"`
	NumericStatus int    `json:"GpoStatus"`
	Status        string
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

// GetGPOFromBytes unmarshalls the incoming byte array containing JSON
// into a GPO structure.
func GetGPOFromBytes(input []byte) (*GPO, error) {
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
func GetGPOFromHost(conn *winrm.Client, name, guid string) (*GPO, error) {
	var cmd string
	if name != "" {
		cmd = getGPOCmdByName(name)
	} else if guid != "" {
		cmd = getGPOCmdByGUID(guid)
	}
	result, err := RunWinRMCommand(conn, []string{cmd}, true)
	if err != nil {
		return nil, err
	}
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("command exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}
	gpo, err := GetGPOFromBytes([]byte(result.Stdout))
	if err != nil {
		return nil, err
	}

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
func (g *GPO) Rename(client *winrm.Client, target string) error {
	if g.ID == "" {
		return fmt.Errorf("gpo guid required")
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("Rename-GPO -Guid %s -TargetName %s", g.ID, g.Name))

	if g.Domain != "" {
		cmds = append(cmds, fmt.Sprintf("-Domain %s", g.Domain))
	}
	cmd := strings.Join(cmds, " ")
	_, err := RunWinRMCommand(client, []string{cmd}, false)
	if err != nil {
		return err
	}
	return nil
}

//ChangeStatus Changes the status of a GPO
func (g *GPO) ChangeStatus(client *winrm.Client, status string) error {
	cmd := fmt.Sprintf(`(%s).GpoStatus = "%s"`, getGPOCmdByGUID(g.ID), status)
	result, err := RunWinRMCommand(client, []string{cmd}, false)
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
func (g *GPO) NewGPO(client *winrm.Client) (string, error) {

	if g.Name == "" {
		return "", fmt.Errorf("gpo name required")
	}
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("New-GPO -Name %s", g.Name))

	if g.Domain != "" {
		cmds = append(cmds, fmt.Sprintf("-Domain %s", g.Domain))
	}

	if g.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Comment '%s'", g.Description))
	}

	result, err := RunWinRMCommand(client, cmds, true)
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
	gpo, err := GetGPOFromBytes([]byte(result.Stdout))
	if err != nil {
		return "", err
	}
	return gpo.ID, nil
}

// DeleteGPO delete the GPO container
func (g *GPO) DeleteGPO(client *winrm.Client) error {
	cmd := fmt.Sprintf("Remove-GPO -Name %s -Domain %s", g.Name, g.Domain)
	_, err := RunWinRMCommand(client, []string{cmd}, false)
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
func (g *GPO) UpdateGPO(client *winrm.Client, d *schema.ResourceData) (string, error) {
	if d.HasChange("name") {
		err := g.Rename(client, SanitiseTFInput(d, "name"))
		if err != nil {
			return "", err
		}
	}

	if d.HasChange("status") {
		err := g.ChangeStatus(client, SanitiseTFInput(d, "status"))
		if err != nil {
			return "", err
		}
	}
	return "", nil
}
