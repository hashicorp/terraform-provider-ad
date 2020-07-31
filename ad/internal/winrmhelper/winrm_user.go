package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
)

// User represents an AD User
type User struct {
	GUID                 string `json:"ObjectGUID"`
	SAMAccountName       string `json:"SamAccountName"`
	PrincipalName        string `json:"UserPrincipalName"`
	DisplayName          string `json:"DisplayName"`
	DistinguishedName    string `json:"DistinguishedName"`
	UserAccountControl   int64  `json:"userAccountControl"`
	Password             string
	Container            string
	Domain               string
	Username             string
	Enabled              bool
	PasswordNeverExpires bool
	CannotChangePassword bool
}

// NewUser creates the user by running the New-ADUser powershell command
func (u *User) NewUser(client *winrm.Client) (string, error) {
	if u.Username == "" {
		return "", fmt.Errorf("user principal name required")
	}

	log.Printf("Adding user with UPN: %q", u.PrincipalName)
	cmds := []string{fmt.Sprintf("New-ADUser -Passthru -Name %q", u.Username)}

	cmds = append(cmds, fmt.Sprintf("-CannotChangePassword $%t", u.CannotChangePassword))
	cmds = append(cmds, fmt.Sprintf("-PasswordNeverExpires $%t", u.PasswordNeverExpires))
	cmds = append(cmds, fmt.Sprintf("-Enabled $%t", u.Enabled))

	if u.SAMAccountName != "" {
		cmds = append(cmds, fmt.Sprintf("-SamAccountName %q", u.SAMAccountName))
	}

	if u.PrincipalName != "" {
		cmds = append(cmds, fmt.Sprintf("-UserPrincipalName %q", u.PrincipalName))
	}

	if u.Password != "" {
		cmds = append(cmds, fmt.Sprintf("-AccountPassword (ConvertTo-SecureString -AsPlainText %q -Force)", u.Password))
	}

	if u.DisplayName != "" {
		cmds = append(cmds, fmt.Sprintf("-DisplayName %q", u.DisplayName))
	}

	if u.Container != "" {
		cmds = append(cmds, fmt.Sprintf("-Path %q", u.Container))
	}

	result, err := RunWinRMCommand(client, cmds, true)
	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		if strings.Contains(result.StdErr, "AlreadyExists") {
			return "", fmt.Errorf("there is another User named %q", u.PrincipalName)
		}
		return "", fmt.Errorf("command New-ADUser exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	user, err := unmarshallUser([]byte(result.Stdout))
	if err != nil {
		return "", fmt.Errorf("error while unmarshalling user json document: %s", err)
	}

	return user.GUID, nil
}

// ModifyUser updates the AD user's details based on what's changed in the resource.
func (u *User) ModifyUser(d *schema.ResourceData, client *winrm.Client) error {
	log.Printf("Modifying user: %q", u.PrincipalName)
	strKeyMap := map[string]string{
		"sam_account_name": "SamAccountName",
		"display_name":     "DisplayName",
		"principal_name":   "UserPrincipalName",
	}

	cmds := []string{fmt.Sprintf("Set-ADUser -Identity %q", u.GUID)}

	for k, param := range strKeyMap {
		if d.HasChange(k) {
			value := d.Get(k).(string)
			cmds = append(cmds, fmt.Sprintf("-%s %q", param, value))
		}
	}

	boolKeyMap := map[string]string{
		"cannot_change_password": "CannotChangePassword",
		"password_never_expires": "PasswordNeverExpires",
		"enabled":                "Enabled",
	}

	for k, param := range boolKeyMap {
		if d.HasChange(k) {
			value := d.Get(k).(bool)
			cmds = append(cmds, fmt.Sprintf("-%s $%t", param, value))
		}
	}

	if d.HasChange("container") {
		cmds = append(cmds, fmt.Sprintf("-Path %q", u.Container))

	}

	if len(cmds) > 1 {
		result, err := RunWinRMCommand(client, cmds, false)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADUser exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("initial_password") {
		cmd := fmt.Sprintf("Set-ADAccountPassword -Identity %q -Reset -NewPassword (ConvertTo-SecureString -AsPlainText %q -Force)", u.GUID, u.Password)
		result, err := RunWinRMCommand(client, []string{cmd}, false)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-AccountPassword exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	return nil
}

//DeleteUser deletes an AD user by calling Remove-ADUser
func (u *User) DeleteUser(client *winrm.Client) error {
	cmd := fmt.Sprintf("Remove-ADUser -Identity %s -Confirm:$false", u.GUID)
	_, err := RunWinRMCommand(client, []string{cmd}, false)
	if err != nil {
		// Check if the resource is already deleted
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	return nil
}

// GetUserFromResource returns a user struct built from Resource data
func GetUserFromResource(d *schema.ResourceData) *User {
	user := User{
		GUID:                 d.Id(),
		SAMAccountName:       SanitiseTFInput(d, "sam_account_name"),
		PrincipalName:        SanitiseTFInput(d, "principal_name"),
		DisplayName:          SanitiseTFInput(d, "display_name"),
		Container:            SanitiseTFInput(d, "container"),
		Password:             SanitiseTFInput(d, "initial_password"),
		Enabled:              d.Get("enabled").(bool),
		PasswordNeverExpires: d.Get("password_never_expires").(bool),
		CannotChangePassword: d.Get("cannot_change_password").(bool),
	}
	if user.PrincipalName != "" {
		tokens := strings.Split(user.PrincipalName, "@")
		user.Username = tokens[0]
		if len(tokens) > 1 {
			user.Domain = tokens[1]
		}
	}

	return &user
}

// GetUserFromHost returns a User struct based on data
// retrieved from the AD Domain Controller.
func GetUserFromHost(client *winrm.Client, guid string) (*User, error) {
	cmd := fmt.Sprintf("Get-ADUser -identity %q -properties *", guid)
	result, err := RunWinRMCommand(client, []string{cmd}, true)
	if err != nil {
		return nil, err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		return nil, fmt.Errorf("command Get-ADUser exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	u, err := unmarshallUser([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling user json document: %s", err)
	}

	return u, nil
}

// unmarshallUser unmarshalls the incoming byte array containing JSON
// into a User structure and populates all fields based on the data
// extracted.
func unmarshallUser(input []byte) (*User, error) {
	var user User
	err := json.Unmarshal(input, &user)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall json document with error %q, document was: %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling json response: %s", err)
	}

	if user.PrincipalName != "" {
		tokens := strings.Split(user.PrincipalName, "@")
		user.Username = tokens[0]
		if len(tokens) > 1 {
			user.Domain = tokens[1]
		}
	}

	commaIdx := strings.Index(user.DistinguishedName, ",")
	user.Container = user.DistinguishedName[commaIdx+1:]

	var accountControlMap = map[string]int64{
		"disabled":               0x00000002,
		"password_never_expires": 0x00010000,
		"cannot_change_password": 0x00000040,
	}

	user.Enabled = !(user.UserAccountControl&accountControlMap["disabled"] != 0)
	user.PasswordNeverExpires = user.UserAccountControl&accountControlMap["password_never_expires"] != 0
	user.CannotChangePassword = user.UserAccountControl&accountControlMap["cannot_change_password"] != 0

	return &user, nil
}
