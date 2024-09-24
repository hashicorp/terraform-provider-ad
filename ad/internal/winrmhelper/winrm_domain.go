package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"
)

// Domain struct represents an AD Domain account object
type Domain struct {
	Name string `json:"Name"`
	GUID string `json:"ObjectGuid"`
	DN   string `json:"DistinguishedName"`
	SID  SID    `json:"DomainSID"`
}

// NewDomainFromResource returns a new Machine struct populated from resource data
func NewDomainFromResource(d *schema.ResourceData) *Domain {
	return &Domain{
		Name: SanitiseTFInput(d, "name"),
		DN:   SanitiseTFInput(d, "dn"),
		GUID: SanitiseTFInput(d, "guid"),
	}
}

// NewDomainFromHost return a new Machine struct populated from data we get
// from the domain controller
func NewDomainFromHost(conf *config.ProviderConf, identity string) (*Domain, error) {
	cmd := fmt.Sprintf("Get-ADDomain -Identity %q", identity)
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
		return nil, fmt.Errorf("winrm execution failure in NewDomainFromHost: %s", err)
	}

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("Get-ADDomain exited with a non zero exit code (%d), stderr: %s", result.ExitCode, result.StdErr)
	}
	domain, err := unmarshallDomain([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("NewDomainFromHost: %s", err)
	}

	return domain, nil
}

func unmarshallDomain(input []byte) (*Domain, error) {
	var domain Domain
	err := json.Unmarshal(input, &domain)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall an ADDomain json document with error %q, document was %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling ADDomain json document: %s", err)
	}
	if domain.GUID == "" {
		return nil, fmt.Errorf("invalid data while unmarshalling Domain data, json doc was: %s", string(input))
	}
	return &domain, nil
}
