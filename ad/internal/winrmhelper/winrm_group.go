package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Group represents an AD Group
type Group struct {
	GUID              string `json:"ObjectGUID"`
	SAMAccountName    string `json:"SamAccountName"`
	Name              string `json:"Name"`
	ScopeNum          int    `json:"GroupScope"`
	CategoryNum       int    `json:"GroupCategory"`
	DistinguishedName string `json:"DistinguishedName"`
	Scope             string
	Category          string
	Container         string
	Description       string
	SID               SID `json:"SID"`
}

// AddGroup creates a new group
func (g *Group) AddGroup(conf *config.ProviderConf) (string, error) {
	log.Printf("[DEBUG] Adding group with name %q", g.Name)
	cmds := []string{fmt.Sprintf("New-ADGroup -Passthru -Name %q -GroupScope %q -GroupCategory %q -Path %q", g.Name, g.Scope, g.Category, g.Container)}

	if g.SAMAccountName != "" {
		cmds = append(cmds, fmt.Sprintf("-SamAccountName %q", g.SAMAccountName))
	}

	if g.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description %q", g.Description))
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
	psCmd := NewPSCommand(cmds, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		if strings.Contains(result.StdErr, "already exists") {
			return "", fmt.Errorf("there is another group named %q", g.Name)
		}
		return "", fmt.Errorf("command New-ADGroup exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	group, err := unmarshallGroup([]byte(result.Stdout))
	if err != nil {
		return "", fmt.Errorf("error while unmarshalling group json document: %s", err)
	}

	return group.GUID, nil
}

// ModifyGroup updates an existing group
func (g *Group) ModifyGroup(d *schema.ResourceData, conf *config.ProviderConf) error {
	KeyMap := map[string]string{
		"sam_account_name": "SamAccountName",
		"scope":            "GroupScope",
		"category":         "GroupCategory",
		"description":      "Description",
	}

	cmds := []string{fmt.Sprintf("Set-ADGroup -Identity %q", g.GUID)}

	for k, param := range KeyMap {
		if d.HasChange(k) {
			value := SanitiseTFInput(d, k)
			if value == "" {
				value = "$null"
			} else {
				value = fmt.Sprintf(`"%s"`, value)
			}
			cmds = append(cmds, fmt.Sprintf(`-%s %s`, param, value))
		}
	}

	if len(cmds) > 1 {
		psOpts := CreatePSCommandOpts{
			JSONOutput:      false,
			ForceArray:      false,
			ExecLocally:     conf.IsConnectionTypeLocal(),
			PassCredentials: conf.IsPassCredentialsEnabled(),
			Username:        conf.Settings.WinRMUsername,
			Password:        conf.Settings.WinRMPassword,
			Server:          conf.IdentifyDomainController(),
		}
		psCmd := NewPSCommand(cmds, psOpts)
		result, err := psCmd.Run(conf)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADGroup exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("name") {
		cmd := fmt.Sprintf("Rename-ADObject -Identity %q -NewName %q", g.GUID, d.Get("name").(string))
		psOpts := CreatePSCommandOpts{
			JSONOutput:      false,
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
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Rename-ADObject exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("container") {
		cmd := fmt.Sprintf("Move-ADObject -Identity %q -TargetPath %q", g.GUID, d.Get("container").(string))
		psOpts := CreatePSCommandOpts{
			JSONOutput:      false,
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
			return fmt.Errorf("winrm execution failure while moving group object: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Move-ADObject exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	return nil
}

// DeleteGroup removes a group
func (g *Group) DeleteGroup(conf *config.ProviderConf) error {
	cmd := fmt.Sprintf("Remove-ADGroup -Identity %s -Confirm:$false", g.GUID)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      false,
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
		// Check if the resource is already deleted
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	} else if result.ExitCode != 0 {
		return fmt.Errorf("while removing group: stderr: %s", result.StdErr)
	}
	return nil
}

// GetGroupFromResource returns a Group struct built from Resource data
func GetGroupFromResource(d *schema.ResourceData) *Group {
	g := Group{
		Name:           SanitiseTFInput(d, "name"),
		SAMAccountName: SanitiseTFInput(d, "sam_account_name"),
		Container:      SanitiseTFInput(d, "container"),
		Scope:          SanitiseTFInput(d, "scope"),
		Category:       SanitiseTFInput(d, "category"),
		GUID:           SanitiseString(d.Id()),
		Description:    SanitiseTFInput(d, "description"),
	}

	return &g
}

// GetGroupFromHost returns a Group struct based on data
// retrieved from the AD Controller.
func GetGroupFromHost(conf *config.ProviderConf, guid string) (*Group, error) {
	cmd := fmt.Sprintf("Get-ADGroup -identity %q -properties *", guid)
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
		return nil, err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		return nil, fmt.Errorf("command Get-ADGroup exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	g, err := unmarshallGroup([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling group json document: %s", err)
	}

	return g, nil
}

// unmarshallGroup unmarshalls the incoming byte array containing JSON
// into a Group structure and populates all fields based on the data
// extracted.
func unmarshallGroup(input []byte) (*Group, error) {
	var g Group
	err := json.Unmarshal(input, &g)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall json document with error %q, document was: %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling json response: %s", err)
	}
	if g.GUID == "" {
		return nil, fmt.Errorf("invalid data while unmarshalling Group data, json doc was: %s", string(input))
	}
	scopes := []string{"domainlocal", "global", "universal"}
	categories := []string{"distribution", "security"}

	g.Scope = scopes[g.ScopeNum]
	g.Category = categories[g.CategoryNum]

	commaIdx := strings.Index(g.DistinguishedName, ",")
	g.Container = g.DistinguishedName[commaIdx+1:]

	return &g, nil
}
