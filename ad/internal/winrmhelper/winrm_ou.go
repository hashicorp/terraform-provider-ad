package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
)

// OrgUnit is a structure used to represent an AD OrganizationalUnit object
type OrgUnit struct {
	Name              string
	Description       string
	Path              string
	Protected         bool `json:"ProtectedFromAccidentalDeletion"`
	DistinguishedName string
	GUID              string `json:"ObjectGuid"`
}

//NewOrgUnitFromResource returns a new OrgUnit struct populated from resource data
func NewOrgUnitFromResource(d *schema.ResourceData) *OrgUnit {
	ou := OrgUnit{
		Description:       SanitiseTFInput(d, "description"),
		Name:              SanitiseTFInput(d, "name"),
		Path:              SanitiseTFInput(d, "path"),
		DistinguishedName: SanitiseTFInput(d, "dn"),
		GUID:              SanitiseTFInput(d, "guid"),
	}
	protected := d.Get("protected").(bool)
	ou.Protected = protected
	return &ou
}

// NewOrgUnitFromHost returns a new OrgUnit struct populated from data we get from
// the domain controller
func NewOrgUnitFromHost(conn *winrm.Client, guid, name, path string, execLocally bool) (*OrgUnit, error) {
	var cmd string
	if guid != "" {
		cmd = fmt.Sprintf("Get-ADObject -Properties * -Identity %q", guid)
	} else if name != "" && path != "" {
		cmd = fmt.Sprintf("Get-ADObject -Properties * -Name %q -Path %q", name, path)
	} else {
		return nil, fmt.Errorf("invalid inputs, dn or a combination of path and name are required")
	}

	result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
	if err != nil {
		return nil, err
	}
	if result.ExitCode != 0 {
		return nil, fmt.Errorf("Get-ADOrganizationalUnit exited with a non-zero exit code %d, stderr :%s", result.ExitCode, result.StdErr)
	}
	ou, err := unmarshallOU([]byte(result.Stdout))
	if err != nil {
		return nil, err
	}
	ou.Path = strings.TrimPrefix(ou.DistinguishedName, fmt.Sprintf("OU=%s,", ou.Name))

	return ou, nil
}

// Create creates a new OU in the AD tree
func (o *OrgUnit) Create(conn *winrm.Client, execLocally bool) (string, error) {

	cmd := "New-ADOrganizationalUnit -Passthru"
	if o.Name == "" {
		return "", fmt.Errorf("missing required attribute name, cannot create OU")
	}
	cmd = fmt.Sprintf("%s -Name %q", cmd, o.Name)

	if o.Description != "" {
		cmd = fmt.Sprintf("%s -Description %q", cmd, o.Description)
	}

	if o.Path != "" {
		cmd = fmt.Sprintf("%s -Path %q", cmd, o.Path)
	}

	cmd = fmt.Sprintf("%s -ProtectedFromAccidentalDeletion:$%t", cmd, o.Protected)

	result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
	if err != nil {
		return "", err
	}
	if result.ExitCode != 0 {
		return "", fmt.Errorf("Get-ADOrganizationalUnit exited with a non-zero exit code %d, stderr :%s", result.ExitCode, result.StdErr)
	}
	ou, err := unmarshallOU([]byte(result.Stdout))
	if err != nil {
		return "", err
	}

	return ou.GUID, nil
}

// Update updates an existing OU in the AD tree
func (o *OrgUnit) Update(conn *winrm.Client, changes map[string]interface{}, execLocally bool) error {
	if o.DistinguishedName == "" {
		return fmt.Errorf("Cannot update OU with name %q, distiguished name is empty", o.Name)
	}
	cmd := fmt.Sprintf("Set-ADOrganizationalUnit -Identity %q", o.DistinguishedName)

	keyMap := map[string]string{
		"display_name": "DisplayName",
		"description":  "Description",
	}

	for k, v := range changes {
		if paramName, ok := keyMap[k]; ok {
			cmd = fmt.Sprintf("%s -%s %q", cmd, paramName, v.(string))
		}
	}

	if cmd != "Set-ADOrganizationalUnit -Identity" {
		result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Set-ADOrganizationalUnit exited with a non-zero exit code %d, stderr :%s", result.ExitCode, result.StdErr)
		}
	}

	if path, ok := changes["path"]; ok {
		var unprotected bool
		if o.Protected == true {
			cmd := fmt.Sprintf("Set-ADOrganizationalUnit -Identity %q -ProtectedFromAccidentalDeletion:$false", o.GUID)
			result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
			if err != nil {
				return fmt.Errorf("winrm execution failure while unprotecting OU object: %s", err)
			}
			if result.ExitCode != 0 {
				return fmt.Errorf("Set-ADOrganizationalUnit exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
			}
			unprotected = true
		}

		cmd := fmt.Sprintf("Move-ADObject -Identity %q -TargetPath %q", o.GUID, path.(string))
		result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
		if err != nil {
			return fmt.Errorf("winrm execution failure while moving OU object: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Move-ADObject exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}

		if unprotected == true {
			cmd := fmt.Sprintf("Set-ADOrganizationalUnit -Identity %q -ProtectedFromAccidentalDeletion:$true", o.GUID)
			result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
			if err != nil {
				return fmt.Errorf("winrm execution failure while protecting OU object: %s", err)
			}
			if result.ExitCode != 0 {
				return fmt.Errorf("Set-ADOrganizationalUnit exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
			}
		}
	}

	if protected, ok := changes["protected"]; ok {
		cmd = fmt.Sprintf("Set-ADObject -Identity %s -ProtectedFromAccidentalDeletion:$%t", o.GUID, protected.(bool))
		result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Set-ADObject exited with a non-zero exit code (%d) while updating OU's protected status, stderr :%s", result.ExitCode, result.StdErr)
		}
	}

	if name, ok := changes["name"]; ok {
		cmd = fmt.Sprintf("Rename-ADObject -Identity %q %q ", o.GUID, name.(string))
		result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Set-ADObject exited with a non-zero exit code (%d) while renaming OU, stderr :%s", result.ExitCode, result.StdErr)
		}

	}
	return nil
}

// Delete deletes an existing OU from an AD tree
func (o *OrgUnit) Delete(conn *winrm.Client, execLocally bool) error {
	if o.DistinguishedName == "" {
		return fmt.Errorf("Cannot remove OU with name %q, distiguished name is empty", o.Name)
	}
	cmd := fmt.Sprintf("Get-ADObject -Properties * -Identity %q | Set-ADObject -ProtectedFromAccidentalDeletion:$false -Passthru | Remove-ADOrganizationalUnit -confirm:$false", o.DistinguishedName)
	result, err := RunWinRMCommand(conn, []string{cmd}, true, false, execLocally)
	if err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("Get-ADObject -Properties * exited with a non-zero exit code %d, stderr :%s", result.ExitCode, result.StdErr)
	}
	return nil
}

func unmarshallOU(input []byte) (*OrgUnit, error) {
	var ou OrgUnit
	err := json.Unmarshal(input, &ou)
	if err != nil {
		log.Printf("[ERROR] Failed to unmarshall json document with error %q, document was: %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling json response: %s", err)
	}
	return &ou, nil

}
