package winrmhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/masterzen/winrm"
)

//GPLink represents an AD Object that links a GPO and another AD object such as a group,
// an OU, or a domain.
type GPLink struct {
	GPOGuid  string `json:"GpoId"`
	Target   string `json:"Target"`
	Enforced bool   `json:"Enforced"`
	Order    int    `json:"Order"`
	Enabled  bool   `json:"Enabled"`
}

//NewGPLink creates a link between a GPO and an AD object
func (g *GPLink) NewGPLink(client *winrm.Client) (string, error) {
	log.Printf("[DEBUG] Creating new user")
	enforced := "No"
	if g.Enforced {
		enforced = "Yes"
	}

	enabled := "No"
	if g.Enabled {
		enabled = "Yes"
	}

	cmds := []string{fmt.Sprintf("New-GPLink -Guid %q -Target %q -LinkEnabled %q -Enforced %q", g.GPOGuid, g.Target, enabled, enforced)}

	if g.Order > 0 {
		cmds = append(cmds, fmt.Sprintf("-Order %d", g.Order))
	}

	result, err := RunWinRMCommand(client, cmds, true)
	if err != nil {
		return "", err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		if strings.Contains(result.StdErr, "is already linked") {
			return "", fmt.Errorf("there is another link between GPO %q and target %q", g.GPOGuid, g.Target)
		}
		return "", fmt.Errorf("command New-GPLink exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	gplink, err := unmarshallNewGPLink([]byte(result.Stdout))
	if err != nil {
		return "", fmt.Errorf("error while unmarshalling gplink json document: %s", err)
	}

	ou, err := NewOrgUnitFromHost(client, gplink.Target, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve details for OU %q: %s", gplink.Target, err)
	}

	id := fmt.Sprintf("%s_%s", gplink.GPOGuid, ou.GUID)

	return id, nil

}

//ModifyGPLink changes a GPO link
func (g *GPLink) ModifyGPLink(client *winrm.Client, changes map[string]interface{}) error {
	cmds := []string{fmt.Sprintf("Set-GPLink -guid %q -target %q", g.GPOGuid, g.Target)}
	keyMap := map[string]string{
		"enforced": "Enforced",
		"enabled":  "LinkEnabled",
	}

	for k, v := range changes {
		if paramName, ok := keyMap[k]; ok {
			value := "No"
			if v.(bool) {
				value = "Yes"
			}
			cmds = append(cmds, fmt.Sprintf("-%s %q", paramName, value))
		}
	}

	if order, ok := changes["order"]; ok {
		cmds = append(cmds, fmt.Sprintf("-Order %d", order.(int)))
	}

	if len(cmds) == 1 {
		return nil
	}
	result, err := RunWinRMCommand(client, cmds, false)
	if err != nil {
		return fmt.Errorf("error while running Set-GPLink: %s", err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("Set-GPLink exited with a non-zero exit code %d, stderr :%s", result.ExitCode, result.StdErr)
	}

	return nil
}

//RemoveGPLink deletes a link between a GPO and an AD object
func (g *GPLink) RemoveGPLink(client *winrm.Client) error {
	cmd := fmt.Sprintf("Remove-GPlink -Guid %q -Target %q", g.GPOGuid, g.Target)
	_, err := RunWinRMCommand(client, []string{cmd}, false)
	if err != nil {
		// Check if the resource is already deleted
		if strings.Contains(err.Error(), "GpoLinkNotFound") || strings.Contains(err.Error(), "GpoWithIdNotFound") || strings.Contains(err.Error(), "There is no such object on the server") {
			return nil
		}
		return err
	}
	return nil
}

//GetGPLinkFromResource returns a GPLink struct populated with data from the
//resource's configuration
func GetGPLinkFromResource(d *schema.ResourceData) *GPLink {
	gplink := GPLink{
		GPOGuid:  SanitiseTFInput(d, "gpo_guid"),
		Target:   SanitiseTFInput(d, "target_dn"),
		Enabled:  d.Get("enabled").(bool),
		Enforced: d.Get("enforced").(bool),
		Order:    d.Get("order").(int),
	}
	return &gplink
}

//GetGPLinkFromHost returns a GPLink struct populated with data retrieved from the
//Domain Controller
func GetGPLinkFromHost(client *winrm.Client, gpoGUID, containerGUID string) (*GPLink, error) {
	cmds := []string{fmt.Sprintf("Get-ADObject -filter '{ObjectGUID -eq %q}' -properties gplink", containerGUID)}
	result, err := RunWinRMCommand(client, cmds, true)
	if err != nil {
		return nil, err
	}

	if result.ExitCode != 0 {
		log.Printf("[DEBUG] stderr: %s\nstdout: %s", result.StdErr, result.Stdout)
		return nil, fmt.Errorf("command New-GPLink exited with a non-zero exit code %d, stderr: %s", result.ExitCode, result.StdErr)
	}

	if result.Stdout == "" {
		return nil, fmt.Errorf("did not find a container with DN %q", containerGUID)
	}

	gplinks, err := getGPLinksFromADObject([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("error while retrieving list of GPOs linked to container %q: %s", containerGUID, err)
	}

	if len(gplinks) == 0 {
		return nil, fmt.Errorf("did not find any GPOs linked to GPO %q", containerGUID)
	}
	gpoFound := false
	gpoOrder := -1
	enforced := false
	enabled := false
	ouDN := ""
	for _, gplink := range gplinks {
		if gplink[0] == gpoGUID {
			gpoFound = true
			order, err := strconv.Atoi(gplink[1])
			if err != nil {
				return nil, fmt.Errorf("GetGPLinkFromHost: error while parsing %q as integer: %s", gplink[1], err)
			}
			gpoOrder = order
			switch gplink[2] {
			case "0":
				enforced = false
				enabled = true
			case "1":
				enforced = false
				enabled = false
			case "2":
				enforced = true
				enabled = true
			case "3":
				enforced = true
				enabled = false
			}
			ouDN = gplink[3]
			break
		}
	}

	if !gpoFound {
		return nil, fmt.Errorf("did not find any GPOs with ID %q attached to container %q", gpoGUID, containerGUID)
	}

	gpo := &GPLink{
		GPOGuid:  gpoGUID,
		Order:    gpoOrder,
		Target:   ouDN,
		Enforced: enforced,
		Enabled:  enabled,
	}

	return gpo, nil
}

func unmarshallNewGPLink(input []byte) (*GPLink, error) {
	var gplink *GPLink
	err := json.Unmarshal(input, &gplink)
	if err != nil {
		log.Printf("[DEBUG] Failed to unmarshall json document with error %q, document was: %s", err, string(input))
		return nil, fmt.Errorf("failed while unmarshalling json response: %s", err)
	}
	return gplink, nil
}

func getGPLinksFromADObject(input []byte) ([][]string, error) {

	type ADObject struct {
		DistinguishedName string `json:"DistinguishedName"`
		GPLink            string `json:"gplink"`
	}

	var ado ADObject
	err := json.Unmarshal(input, &ado)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling Get-ADObject response: %s", err)
	}

	out := [][]string{}
	gpLinks := strings.Split(ado.GPLink, "[")
	re := regexp.MustCompile("{([\\w-]+)}[\\w,=-]+;([0-9])")
	for idx, gpLink := range gpLinks {
		gpoGUIDs := re.FindAllStringSubmatch(gpLink, -1)
		if gpoGUIDs != nil && len(gpoGUIDs) == 1 && len(gpoGUIDs[0]) == 3 {
			// gpoGUIDs has three elements. First is the whole matched string,
			// second is the GPO GUID and third is the gpLinkOptions field
			out = append(out, []string{gpoGUIDs[0][1], fmt.Sprintf("%d", idx), gpoGUIDs[0][2], ado.DistinguishedName})
		}
	}
	return out, nil
}
