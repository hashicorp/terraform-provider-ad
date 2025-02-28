package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Computer struct represents an AD Computer account object
type Computer struct {
	Name           string
	GUID           string `json:"ObjectGuid"`
	DN             string `json:"DistinguishedName"`
	Description    string
	SAMAccountName string `json:"SamAccountName"`
	Path           string
	SID            SID `json:"SID"`
}

// NewComputerFromResource returns a new Machine struct populated from resource data
func NewComputerFromResource(d *schema.ResourceData) *Computer {
	return &Computer{
		Name:           SanitiseTFInput(d, "name"),
		DN:             SanitiseTFInput(d, "dn"),
		Description:    SanitiseTFInput(d, "description"),
		GUID:           SanitiseTFInput(d, "guid"),
		SAMAccountName: SanitiseTFInput(d, "pre2kname"),
		Path:           SanitiseTFInput(d, "container"),
	}
}

// NewComputerFromHost return a new Machine struct populated from data we get
// from the domain controller
func NewComputerFromHost(conf *config.ProviderConf, identity string) (*Computer, error) {
	cmd := fmt.Sprintf("Get-ADComputer -Identity %q -Properties *", identity)
	conn, err := conf.AcquireWinRMClient()
	if err != nil {
		return nil, fmt.Errorf("while acquiring winrm client: %s", err)
	}
	defer conf.ReleaseWinRMClient(conn)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      true,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.IdentifyDomainController(),
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return nil, fmt.Errorf("winrm execution failure in NewComputerFromHost: %s", err)
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("Get-ADComputer exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
	}
	computer, err := unmarshallComputer([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("NewComputerFromHost: %s", err)
	}
	computer.Path = strings.TrimPrefix(computer.DN, fmt.Sprintf("CN=%s,", computer.Name))

	return computer, nil
}

// Create creates a new Computer object in the AD tree
func (m *Computer) Create(conf *config.ProviderConf) (string, error) {
	if m.Name == "" {
		return "", fmt.Errorf("Computer.Create: missing name variable")
	}
	cmd := fmt.Sprintf("New-ADComputer -Passthru -Name %q", m.Name)

	if m.SAMAccountName != "" {
		cmd = fmt.Sprintf("%s -SamAccountName %q", cmd, m.SAMAccountName)
	}

	if m.Path != "" {
		cmd = fmt.Sprintf("%s -Path %q", cmd, m.Path)
	}

	if m.Description != "" {
		cmd = fmt.Sprintf("%s -Description %q", cmd, m.Description)
	}

	psOpts := CreatePSCommandOpts{
		JSONOutput:      true,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.IdentifyDomainController(),
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return "", fmt.Errorf("winrm execution failure while creating computer object: %s", err)
	}

	if result.ExitCode != 0 {
		return "", fmt.Errorf("New-ADComputer exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
	}
	computer, err := unmarshallComputer([]byte(result.Stdout))
	if err != nil {
		return "", fmt.Errorf("Computer.Create: %s", err)
	}

	return computer.GUID, nil
}

// Update updates an existing Computer objects in the AD tree
func (m *Computer) Update(conf *config.ProviderConf, changes map[string]interface{}) error {
	if m.GUID == "" {
		return fmt.Errorf("cannot update computer object with name %q, guid is not set", m.Name)
	}

	if path, ok := changes["container"]; ok {
		cmd := fmt.Sprintf("Move-AdObject -Identity %q -TargetPath %q", m.GUID, path.(string))
		conn, err := conf.AcquireWinRMClient()
		if err != nil {
			return fmt.Errorf("while acquiring winrm client: %s", err)
		}
		defer conf.ReleaseWinRMClient(conn)
		psOpts := CreatePSCommandOpts{
			JSONOutput:      true,
			ForceArray:      false,
			ExecLocally:     conf.IsConnectionTypeLocal(),
			PassCredentials: conf.IsPassCredentialsEnabled(),
			Username:        conf.Settings.WinRMUsername,
			Password:        conf.Settings.WinRMPassword,
			Server:          conf.IdentifyDomainController(),
		}
		psCmd := NewPSCommand([]string{cmd}, psOpts)
		result, err := psCmd.Run(conf)
		if err != nil {
			return fmt.Errorf("winrm execution failure while moving computer object: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Move-ADObject exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if description, ok := changes["description"]; ok {
		if description == "" {
			description = "$null"
		} else {
			description = fmt.Sprintf("%q", description)
		}
		cmd := fmt.Sprintf("Set-ADComputer -Identity %q -Description %s", m.GUID, description)
		conn, err := conf.AcquireWinRMClient()
		if err != nil {
			return fmt.Errorf("while acquiring winrm client: %s", err)
		}
		defer conf.ReleaseWinRMClient(conn)
		psOpts := CreatePSCommandOpts{
			JSONOutput:      true,
			ForceArray:      false,
			ExecLocally:     conf.IsConnectionTypeLocal(),
			PassCredentials: conf.IsPassCredentialsEnabled(),
			Username:        conf.Settings.WinRMUsername,
			Password:        conf.Settings.WinRMPassword,
			Server:          conf.IdentifyDomainController(),
		}
		psCmd := NewPSCommand([]string{cmd}, psOpts)
		result, err := psCmd.Run(conf)
		if err != nil {
			return fmt.Errorf("winrm execution failure while modifying computer description: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Set-ADComputer exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	return nil
}

// Delete deletes an existing Computer objects from the AD tree
func (m *Computer) Delete(conf *config.ProviderConf) error {
	cmd := fmt.Sprintf("Remove-ADObject -Confirm:$false -Recursive -Identity %q", m.GUID)
	conn, err := conf.AcquireWinRMClient()
	if err != nil {
		return fmt.Errorf("while acquiring winrm client: %s", err)
	}
	defer conf.ReleaseWinRMClient(conn)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      true,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.IdentifyDomainController(),
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return fmt.Errorf("winrm execution failure while removing computer object: %s", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("Remove-ADComputer exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
	}
	return nil
}

func unmarshallComputer(input []byte) (*Computer, error) {
	var computer Computer
	err := json.Unmarshal(input, &computer)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall an ADComputer json document with error %q, document was %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling ADComputer json document: %s", err)
	}
	if computer.GUID == "" {
		return nil, fmt.Errorf("invalid data while unmarshalling Computer data, json doc was: %s", string(input))
	}
	return &computer, nil
}
