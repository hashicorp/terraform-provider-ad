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
	GUID                   string `json:"ObjectGUID"`
	SAMAccountName         string `json:"SamAccountName"`
	PrincipalName          string `json:"UserPrincipalName"`
	City                   string
	Company                string
	Country                string
	Department             string
	Description            string
	DisplayName            string `json:"DisplayName"`
	DistinguishedName      string `json:"DistinguishedName"`
	Division               string
	EmailAddress           string
	EmployeeID             string
	EmployeeNumber         string
	Enabled                bool
	Fax                    string
	GivenName              string
	HomeDirectory          string
	HomeDrive              string
	HomePhone              string
	HomePage               string
	Initials               string
	MobilePhone            string
	Office                 string
	OfficePhone            string
	Organization           string
	OtherName              string
	POBox                  string
	PostalCode             string
	SmartcardLogonRequired bool
	State                  string
	StreetAddress          string
	Surname                string
	Title                  string
	TrustedForDelegation   bool
	UserAccountControl     int64 `json:"userAccountControl"`
	Password               string
	Container              string
	Domain                 string
	Username               string
	PasswordNeverExpires   bool
	CannotChangePassword   bool
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
	cmds = append(cmds, fmt.Sprintf("-SmartcardLogonRequired $%t", u.SmartcardLogonRequired))
	cmds = append(cmds, fmt.Sprintf("-TrustedForDelegation $%t", u.TrustedForDelegation))

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

	if u.City != "" {
		cmds = append(cmds, fmt.Sprintf("-City %q", u.City))
	}

	if u.Company != "" {
		cmds = append(cmds, fmt.Sprintf("-Company %q", u.Company))
	}

	if u.Country != "" {
		cmds = append(cmds, fmt.Sprintf("-Country %q", u.Country))
	}

	if u.Department != "" {
		cmds = append(cmds, fmt.Sprintf("-Department %q", u.Department))
	}

	if u.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description %q", u.Description))
	}

	if u.Division != "" {
		cmds = append(cmds, fmt.Sprintf("-Division %q", u.Division))
	}

	if u.EmailAddress != "" {
		cmds = append(cmds, fmt.Sprintf("-EmailAddress %q", u.EmailAddress))
	}

	if u.EmployeeID != "" {
		cmds = append(cmds, fmt.Sprintf("-EmployeeID %q", u.EmployeeID))
	}

	if u.EmployeeNumber != "" {
		cmds = append(cmds, fmt.Sprintf("-EmployeeNumber %q", u.EmployeeNumber))
	}

	if u.Fax != "" {
		cmds = append(cmds, fmt.Sprintf("-Fax %q", u.Fax))
	}

	if u.GivenName != "" {
		cmds = append(cmds, fmt.Sprintf("-GivenName %q", u.GivenName))
	}

	if u.HomeDirectory != "" {
		cmds = append(cmds, fmt.Sprintf("-HomeDirectory %q", u.HomeDirectory))
	}

	if u.HomeDrive != "" {
		cmds = append(cmds, fmt.Sprintf("-HomeDrive %q", u.HomeDrive))
	}

	if u.HomePhone != "" {
		cmds = append(cmds, fmt.Sprintf("-HomePhone %q", u.HomePhone))
	}

	if u.HomePage != "" {
		cmds = append(cmds, fmt.Sprintf("-HomePage %q", u.HomePage))
	}

	if u.Initials != "" {
		cmds = append(cmds, fmt.Sprintf("-Initials %q", u.Initials))
	}

	if u.MobilePhone != "" {
		cmds = append(cmds, fmt.Sprintf("-MobilePhone %q", u.MobilePhone))
	}

	if u.Office != "" {
		cmds = append(cmds, fmt.Sprintf("-Office %q", u.Office))
	}

	if u.OfficePhone != "" {
		cmds = append(cmds, fmt.Sprintf("-OfficePhone %q", u.OfficePhone))
	}

	if u.Organization != "" {
		cmds = append(cmds, fmt.Sprintf("-Organization %q", u.Organization))
	}

	if u.OtherName != "" {
		cmds = append(cmds, fmt.Sprintf("-OtherName %q", u.OtherName))
	}

	if u.POBox != "" {
		cmds = append(cmds, fmt.Sprintf("-POBox %q", u.POBox))
	}

	if u.PostalCode != "" {
		cmds = append(cmds, fmt.Sprintf("-PostalCode %q", u.PostalCode))
	}

	if u.State != "" {
		cmds = append(cmds, fmt.Sprintf("-State %q", u.State))
	}

	if u.StreetAddress != "" {
		cmds = append(cmds, fmt.Sprintf("-StreetAddress %q", u.StreetAddress))
	}

	if u.Surname != "" {
		cmds = append(cmds, fmt.Sprintf("-Surname %q", u.Surname))
	}

	if u.Title != "" {
		cmds = append(cmds, fmt.Sprintf("-Title %q", u.Title))
	}

	result, err := RunWinRMCommand(client, cmds, true, false)
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
		"city":             "City",
		"company":          "Company",
		"country":          "Country",
		"department":       "Department",
		"description":      "Description",
		"division":         "Division",
		"email_address":    "EmailAddress",
		"employee_id":      "EmployeeID",
		"employee_number":  "EmployeeNumber",
		"fax":              "Fax",
		"given_name":       "GivenName",
		"home_directory":   "HomeDirectory",
		"home_drive":       "HomeDrive",
		"home_phone":       "HomePhone",
		"home_page":        "HomePage",
		"initials":         "Initials",
		"mobile_phone":     "MobilePhone",
		"office":           "Office",
		"office_phone":     "OfficePhone",
		"organization":     "Organization",
		"other_name":       "OtherName",
		"po_box":           "POBox",
		"postal_code":      "PostalCode",
		"state":            "State",
		"street_address":   "StreetAddress",
		"surname":          "Surname",
		"title":            "Title",
	}

	cmds := []string{fmt.Sprintf("Set-ADUser -Identity %q", u.GUID)}

	for k, param := range strKeyMap {
		if d.HasChange(k) {
			value := d.Get(k).(string)
			cmds = append(cmds, fmt.Sprintf("-%s %q", param, value))
		}
	}

	boolKeyMap := map[string]string{
		"cannot_change_password":    "CannotChangePassword",
		"password_never_expires":    "PasswordNeverExpires",
		"enabled":                   "Enabled",
		"smart_card_logon_required": "SmartcardLogonRequired",
		"trusted_for_delegation":    "TrustedForDelegation",
	}

	for k, param := range boolKeyMap {
		if d.HasChange(k) {
			value := d.Get(k).(bool)
			cmds = append(cmds, fmt.Sprintf("-%s $%t", param, value))
		}
	}

	if len(cmds) > 1 {
		result, err := RunWinRMCommand(client, cmds, false, false)
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
		result, err := RunWinRMCommand(client, []string{cmd}, false, false)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-AccountPassword exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("container") {
		path := d.Get("container").(string)
		cmd := fmt.Sprintf("Move-AdObject -Identity %q -TargetPath %q", u.GUID, path)
		result, err := RunWinRMCommand(client, []string{cmd}, true, false)
		if err != nil {
			return fmt.Errorf("winrm execution failure while moving user object: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Move-ADObject exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	return nil
}

//DeleteUser deletes an AD user by calling Remove-ADUser
func (u *User) DeleteUser(client *winrm.Client) error {
	cmd := fmt.Sprintf("Remove-ADUser -Identity %s -Confirm:$false", u.GUID)
	_, err := RunWinRMCommand(client, []string{cmd}, false, false)
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
		GUID:                   d.Id(),
		SAMAccountName:         SanitiseTFInput(d, "sam_account_name"),
		PrincipalName:          SanitiseTFInput(d, "principal_name"),
		DisplayName:            SanitiseTFInput(d, "display_name"),
		Container:              SanitiseTFInput(d, "container"),
		Password:               SanitiseTFInput(d, "initial_password"),
		Enabled:                d.Get("enabled").(bool),
		PasswordNeverExpires:   d.Get("password_never_expires").(bool),
		CannotChangePassword:   d.Get("cannot_change_password").(bool),
		City:                   SanitiseTFInput(d, "city"),
		Company:                SanitiseTFInput(d, "company"),
		Country:                SanitiseTFInput(d, "country"),
		Department:             SanitiseTFInput(d, "department"),
		Description:            SanitiseTFInput(d, "description"),
		Division:               SanitiseTFInput(d, "division"),
		EmailAddress:           SanitiseTFInput(d, "email_address"),
		EmployeeID:             SanitiseTFInput(d, "employee_id"),
		EmployeeNumber:         SanitiseTFInput(d, "employee_number"),
		Fax:                    SanitiseTFInput(d, "fax"),
		GivenName:              SanitiseTFInput(d, "given_name"),
		HomeDirectory:          SanitiseTFInput(d, "home_directory"),
		HomeDrive:              SanitiseTFInput(d, "home_drive"),
		HomePhone:              SanitiseTFInput(d, "home_phone"),
		HomePage:               SanitiseTFInput(d, "home_page"),
		Initials:               SanitiseTFInput(d, "initials"),
		MobilePhone:            SanitiseTFInput(d, "mobile_phone"),
		Office:                 SanitiseTFInput(d, "office"),
		OfficePhone:            SanitiseTFInput(d, "office_phone"),
		Organization:           SanitiseTFInput(d, "organization"),
		OtherName:              SanitiseTFInput(d, "other_name"),
		POBox:                  SanitiseTFInput(d, "po_box"),
		PostalCode:             SanitiseTFInput(d, "postal_code"),
		SmartcardLogonRequired: d.Get("smart_card_logon_required").(bool),
		State:                  SanitiseTFInput(d, "state"),
		StreetAddress:          SanitiseTFInput(d, "street_address"),
		Surname:                SanitiseTFInput(d, "surname"),
		Title:                  SanitiseTFInput(d, "title"),
		TrustedForDelegation:   d.Get("trusted_for_delegation").(bool),
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
	result, err := RunWinRMCommand(client, []string{cmd}, true, false)
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
