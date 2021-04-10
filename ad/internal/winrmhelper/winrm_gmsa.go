package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
)

// Gmsa represents an AD Gmsa
type Gmsa struct {
	AllowReversiblePasswordEncryption          bool
	CanonicalName                              string
	CN                                         string
	Container                                  string
	Created                                    string
	Delegated                                  bool `json:"AccountNotDelegated"`
	Description                                string
	DisplayName                                string
	DistinguishedName                          string
	DNSHostName                                string
	DoesNotRequirePreAuth                      bool
	Enabled                                    bool
	ExpirationString                           string `json:"AccountExpirationDate"`
	Expiration                                 string
	GUID                                       string `json:"ObjectGUID"`
	HomedirRequired                            bool
	HomePage                                   string
	KerberosEncryptionTypeNum                  []int `json:"KerberosEncryptionType"`
	KerberosEncryptionType                     []string
	LastLogonDate                              string
	LockedOut                                  bool
	logonCount                                 int
	ManagedPasswordIntervalInDays              int `json:"msDS-ManagedPasswordInterval"`
	Name                                       string
	PrimaryGroup                               string
	PrincipalsAllowedToDelegateToAccount       []string
	PrincipalsAllowedToRetrieveManagedPassword []string
	SAMAccountName                             string
	ServicePrincipalNames                      []string
	SID                                        SID `json:"SID"`
	TrustedForDelegation                       bool
	UserAccountControl                         int64 `json:"userAccountControl"`
}

// NewGmsa creates the Gmsa by running the New-ADServiceAccount powershell command
func (g *Gmsa) NewGmsa(client *winrm.Client, execLocally bool) (string, error) {
	if g.Name == "" || g.DNSHostName == "" {
		return "", fmt.Errorf("Gmsa name and dnshostname are required !")
	}

	log.Printf("Adding gmsa with name: %q", g.Name)
	cmds := []string{fmt.Sprintf("New-ADServiceAccount -Passthru -Name %q -DNSHostName %q", g.Name, g.DNSHostName)}

	if g.Container != "" {
		cmds = append(cmds, fmt.Sprintf("-Path %q", g.Container))
	}

	if g.Delegated == true {
		cmds = append(cmds, fmt.Sprintf("-AccountNotDelegated $%t", g.Delegated))
	}

	if g.Description != "" {
		cmds = append(cmds, fmt.Sprintf("-Description %q", g.Description))
	}

	if g.DisplayName != "" {
		cmds = append(cmds, fmt.Sprintf("-DisplayName %q", g.DisplayName))
	}

	cmds = append(cmds, fmt.Sprintf("-Enabled $%t", g.Enabled))

	if g.Expiration != "" {
		cmds = append(cmds, fmt.Sprintf("-AccountExpirationDate %q", g.Expiration))
	}

	if g.HomePage != "" {
		cmds = append(cmds, fmt.Sprintf("-HomePage %q", g.HomePage))
	}

	if len(g.KerberosEncryptionType) > 0 {
		cmds = append(cmds, fmt.Sprintf("-KerberosEncryptionType %q", strings.Join(g.KerberosEncryptionType, ",")))
	}

	if g.ManagedPasswordIntervalInDays != 0 {
		cmds = append(cmds, fmt.Sprintf("-ManagedPasswordIntervalInDays %d", g.ManagedPasswordIntervalInDays))
	}

	if len(g.PrincipalsAllowedToDelegateToAccount) > 0 {
		cmds = append(cmds, fmt.Sprintf("-PrincipalsAllowedToDelegateToAccount %q", strings.Join(g.PrincipalsAllowedToDelegateToAccount, ",")))
	}

	if len(g.PrincipalsAllowedToRetrieveManagedPassword) > 0 {
		cmds = append(cmds, fmt.Sprintf("-PrincipalsAllowedToRetrieveManagedPassword %q", strings.Join(g.PrincipalsAllowedToRetrieveManagedPassword, ",")))
	}

	if g.SAMAccountName != "" {
		cmds = append(cmds, fmt.Sprintf("-SamAccountName %q", g.SAMAccountName))
	} else {
		cmds = append(cmds, fmt.Sprintf("-SamAccountName %q", g.Name))
	}

	if len(g.ServicePrincipalNames) > 0 {
		cmds = append(cmds, fmt.Sprintf("-ServicePrincipalNames %q", strings.Join(g.ServicePrincipalNames, ",")))
	}

	cmds = append(cmds, fmt.Sprintf("-TrustedForDelegation $%t", g.TrustedForDelegation))

	result, err := RunWinRMCommand(client, cmds, true, false, execLocally)
	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		if strings.Contains(result.StdErr, "AlreadyExists") {
			return "", fmt.Errorf("there is another gmsa named %q", g.Name)
		}
		return "", fmt.Errorf("command New-ADgmsa exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	gmsa, err := unmarshallGmsa([]byte(result.Stdout))

	if err != nil {
		return "", fmt.Errorf("error while unmarshalling gmsa json document: %s", err)
	}

	return gmsa.GUID, nil
}

// ModifyGmsa updates the AD gmsa's details based on what's changed in the resource.
func (g *Gmsa) ModifyGmsa(d *schema.ResourceData, client *winrm.Client, execLocally bool) error {
	log.Printf("Modifying gmsa: %q", g.Name)
	strKeyMap := map[string]string{
		"expiration":       "AccountExpirationDate",
		"sam_account_name": "SamAccountName",
		"display_name":     "DisplayName",
		"description":      "Description",
		"dns_host_name":    "DNSHostName",
		"home_page":        "HomePage",
		"name":             "Name",
	}

	cmds := []string{fmt.Sprintf("Set-ADServiceAccount -Identity %q", g.GUID)}

	for k, param := range strKeyMap {
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

	boolKeyMap := map[string]string{
		"delegated":              "AccountNotDelegated ",
		"enabled":                "Enabled",
		"trusted_for_delegation": "TrustedForDelegation",
	}

	for k, param := range boolKeyMap {
		if d.HasChange(k) {
			value := d.Get(k).(bool)
			cmds = append(cmds, fmt.Sprintf("-%s $%t", param, value))
		}
	}

	if len(cmds) > 1 {
		result, err := RunWinRMCommand(client, cmds, false, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADServiceAccount exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("service_principal_names") {
		sp := []string{}
		sprinc := d.Get("service_principal_names").(*schema.Set)
		for _, s := range sprinc.List() {
			if s == "" {
				continue
			}
			sp = append(sp, s.(string))
		}
		cmds := []string{fmt.Sprintf("Set-ADServiceAccount -Identity %q -ServicePrincipalNames $null", g.GUID)}
		if len(sp) > 0 {
			spnlist := strings.Join(sp, ",")
			log.Printf("[DEBUG] SPN list: %s", spnlist)
			cmds = append(cmds, fmt.Sprintf("; Set-ADServiceAccount -Identity %q -ServicePrincipalNames @{Add=%q}", g.GUID, spnlist))
		}

		result, err := RunWinRMCommand(client, cmds, false, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADServiceAccount -ServicePrincipalNames exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("principals_allowed_to_delegate_to_account") {
		del := []string{}
		delegate := d.Get("principals_allowed_to_delegate_to_account").(*schema.Set)
		for _, d := range delegate.List() {
			if d == "" {
				continue
			}
			del = append(del, fmt.Sprintf("%q", d.(string)))
		}

		cmds := []string{fmt.Sprintf("Set-ADServiceAccount -Identity %q -PrincipalsAllowedToDelegateToAccount $null", g.GUID)}
		if len(del) > 0 {
			princ_del := strings.Join(del, ",")
			log.Printf("[DEBUG] Principal list: %s", princ_del)
			cmds = append(cmds, fmt.Sprintf(" ; Set-ADServiceAccount -Identity %q -PrincipalsAllowedToDelegateToAccount %s", g.GUID, princ_del))
		}
		result, err := RunWinRMCommand(client, cmds, false, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADServiceAccount -PrincipalsAllowedToDelegateToAccount exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("principals_allowed_to_retrieve_managed_password") {
		pass := []string{}
		passwords := d.Get("principals_allowed_to_retrieve_managed_password").(*schema.Set)
		for _, p := range passwords.List() {
			if p == "" {
				continue
			}
			pass = append(pass, fmt.Sprintf("%q", p.(string)))
		}
		cmds := []string{fmt.Sprintf("Set-ADServiceAccount -Identity %q -PrincipalsAllowedToRetrieveManagedPassword $null", g.GUID)}
		if len(pass) > 0 {
			princ_pass := strings.Join(pass, ",")
			log.Printf("[DEBUG] Principal list: %s", princ_pass)
			cmds = append(cmds, fmt.Sprintf(" ; Set-ADServiceAccount -Identity %q -PrincipalsAllowedToRetrieveManagedPassword %s", g.GUID, princ_pass))
		}
		result, err := RunWinRMCommand(client, cmds, false, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADServiceAccount -PrincipalsAllowedToRetrieveManagedPassword exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("kerberos_encryption_type") {
		krb := []string{}
		krbenc := d.Get("kerberos_encryption_type").(*schema.Set)
		for _, k := range krbenc.List() {
			if k == "" {
				continue
			}
			krb = append(krb, fmt.Sprintf("%q", k.(string)))
		}
		kerb_enc := strings.Join(krb, ",")

		log.Printf("[DEBUG] Kerberos encryption list: %s", kerb_enc)

		cmd := fmt.Sprintf("Set-ADServiceAccount -Identity %q -KerberosEncryptionType %s", g.GUID, kerb_enc)
		result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally)
		if err != nil {
			return err
		}
		if result.ExitCode != 0 {
			log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
			return fmt.Errorf("command Set-ADServiceAccount -KerberosEncryptionType exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	if d.HasChange("container") {
		path := d.Get("container").(string)
		cmd := fmt.Sprintf("Move-ADObject -Identity %q -TargetPath %q", g.GUID, path)
		result, err := RunWinRMCommand(client, []string{cmd}, true, false, execLocally)
		if err != nil {
			return fmt.Errorf("winrm execution failure while moving gmsa object: %s", err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("Move-ADObject exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
		}
	}

	return nil
}

// DeleteGmsa deletes an AD gmsa by calling Remove-ADServiceAccount
func (g *Gmsa) DeleteGmsa(client *winrm.Client, execLocally bool) error {
	cmd := fmt.Sprintf("Remove-ADServiceAccount -Identity %s -Confirm:$false", g.GUID)
	_, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally)
	if err != nil {
		// Check if the resource is already deleted
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return err
	}
	return nil
}

// GetGmsaFromResource returns a gmsa struct built from Resource data
func GetGmsaFromResource(d *schema.ResourceData) *Gmsa {
	gmsa := Gmsa{
		Container:                     SanitiseTFInput(d, "container"),
		Delegated:                     d.Get("delegated").(bool),
		Description:                   SanitiseTFInput(d, "description"),
		DisplayName:                   SanitiseTFInput(d, "display_name"),
		DNSHostName:                   SanitiseTFInput(d, "dns_host_name"),
		Enabled:                       d.Get("enabled").(bool),
		Expiration:                    SanitiseTFInput(d, "expiration"),
		GUID:                          d.Id(),
		HomePage:                      SanitiseTFInput(d, "home_page"),
		ManagedPasswordIntervalInDays: d.Get("managed_password_interval_in_days").(int),
		Name:                          SanitiseTFInput(d, "name"),

		TrustedForDelegation: d.Get("trusted_for_delegation").(bool),
	}

	// delegate
	del := []string{}
	delegate := d.Get("principals_allowed_to_delegate_to_account").(*schema.Set)
	for _, d := range delegate.List() {
		if d == "" {
			continue
		}
		del = append(del, d.(string))
	}
	if del != nil {
		gmsa.PrincipalsAllowedToDelegateToAccount = del
	}

	// principal(s) allowed to retreieve password
	pass := []string{}
	passwords := d.Get("principals_allowed_to_retrieve_managed_password").(*schema.Set)
	for _, p := range passwords.List() {
		if p == "" {
			continue
		}
		pass = append(pass, p.(string))
	}

	if pass != nil {
		gmsa.PrincipalsAllowedToRetrieveManagedPassword = pass
	}

	// kerberos encryption types
	krb := []string{}
	krbenc := d.Get("kerberos_encryption_type").(*schema.Set)
	for _, k := range krbenc.List() {
		if k == "" {
			continue
		}
		krb = append(krb, strings.ToUpper(k.(string)))
	}

	if krb != nil {
		gmsa.KerberosEncryptionType = krb
	}

	return &gmsa
}

// GetGmsaFromHost returns a gmsa struct based on data
// retrieved from the AD Domain Controller.
func GetGmsaFromHost(client *winrm.Client, guid string, execLocally bool) (*Gmsa, error) {
	cmd := fmt.Sprintf("Get-ADServiceAccount -identity %q -properties *", guid)
	result, err := RunWinRMCommand(client, []string{cmd}, true, false, execLocally)
	if err != nil {
		return nil, err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		return nil, fmt.Errorf("command Get-ADgmsa exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	g, err := unmarshallGmsa([]byte(result.Stdout))

	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling gmsa json document: %s", err)
	}

	return g, nil
}

// unmarshallGmsa unmarshalls the incoming byte array containing JSON
// into a Gmsa structure and populates all fields based on the data
// extracted.
func unmarshallGmsa(input []byte) (*Gmsa, error) {
	var gmsa Gmsa
	err := json.Unmarshal(input, &gmsa)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall a ADGmsa json document with error %q, document was %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling ADGmsa json document: %s", err)
	}

	commaIdx := strings.Index(gmsa.DistinguishedName, ",")
	gmsa.Container = gmsa.DistinguishedName[commaIdx+1:]

	if gmsa.ExpirationString != "" {
		log.Printf("[DEBUG] unmarshall :: converting expiration date to proper format (current value: %s)", gmsa.ExpirationString)
		var regdate = regexp.MustCompile(`^\/Date\((.+)\)\/$`)
		// extract unixtime date
		match := regdate.FindStringSubmatch(gmsa.ExpirationString)

		if len(match) == 0 {
			return nil, fmt.Errorf("Failed to unmarshall a ADGmsa json, expiration date format is not matching regex.")
		}

		log.Printf("[DEBUG] unmarshall :: unixtimestamp extracted from AccountExpirationDate attribute: %q", match[1])
		// convert string date to int64
		log.Printf("[DEBUG] unmarshall :: converting unixtimestamp to %s int64", match[1])
		n, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed while unmarshalling ADGmsa json document: %s", err)
		}
		// convert unix_timestamp to RFC3339
		log.Printf("[DEBUG] unmarshall :: converting unixtimestamp int64 to RFC3339")
		t := time.Unix(0, n*int64(time.Millisecond))
		tst := t.Format(time.RFC3339)
		log.Printf("[DEBUG] unmarshall :: converted unixtimestamp to RFC3339 : %s", tst)
		gmsa.Expiration = tst
	}

	if gmsa.KerberosEncryptionTypeNum != nil {
		log.Printf("[DEBUG] Converting Keberos encryption type number to string slice")
		for _, k := range gmsa.KerberosEncryptionTypeNum {
			log.Printf("[DEBUG] Keberos encryption type number : %d", k)
			var krblist []string
			switch k {
			case 4:
				krblist = []string{"RC4"}
			case 8:
				krblist = []string{"AES128"}
			case 12:
				krblist = []string{"RC4", "AES128"}
			case 16:
				krblist = []string{"AES256"}
			case 20:
				krblist = []string{"RC4", "AES256"}
			case 24:
				krblist = []string{"AES128", "AES256"}
			case 28:
				krblist = []string{"RC4", "AES128", "AES256"}
			}
			log.Printf("[DEBUG] Keberos slice list : %q", krblist)
			gmsa.KerberosEncryptionType = krblist
		}
	}

	return &gmsa, nil
}
