package winrmhelper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GroupMembership struct {
	GroupGUID    string
	GroupMembers []*GroupMember
}

type GroupMember struct {
	SamAccountName string `json:"SamAccountName"`
	DN             string `json:"DistinguishedName"`
	GUID           string `json:"ObjectGUID"`
	Name           string `json:"Name"`
}

func groupExistsInList(g *GroupMember, memberList []*GroupMember) bool {
	for _, item := range memberList {
		if g.GUID == item.GUID {
			return true
		}
	}
	return false
}

func diffGroupMemberLists(expectedMembers, existingMembers []*GroupMember) ([]*GroupMember, []*GroupMember) {
	var toAdd, toRemove []*GroupMember
	for _, member := range expectedMembers {
		if !groupExistsInList(member, existingMembers) {
			toAdd = append(toAdd, member)
		}
	}

	for _, member := range existingMembers {
		if !groupExistsInList(member, expectedMembers) {
			toRemove = append(toRemove, member)
		}
	}

	return toAdd, toRemove
}

func unmarshalGroupMembership(input []byte) ([]*GroupMember, error) {
	var gm []*GroupMember
	err := json.Unmarshal(input, &gm)
	if err != nil {
		return nil, err
	}
	if len(gm) > 0 && gm[0].GUID == "" {
		return nil, fmt.Errorf("invalid data while unmarshalling group membership data, json doc was: %s", string(input))
	}
	return gm, nil
}

func getMembershipList(g []*GroupMember) string {
	out := []string{}
	for _, member := range g {
		out = append(out, member.GUID)
	}

	return strings.Join(out, ",")
}

func (g *GroupMembership) getGroupMembers(conf *config.ProviderConf) ([]*GroupMember, error) {
	cmd := fmt.Sprintf("Get-ADGroupMember -Identity %q", g.GroupGUID)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      true,
		ForceArray:      true,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return nil, fmt.Errorf("while running Get-ADGroupMember: %s", err)
	} else if result.ExitCode != 0 {
		return nil, fmt.Errorf("command Get-ADGroupMember exited with a non-zero exit code(%d), stderr: %s, stdout: %s", result.ExitCode, result.StdErr, result.Stdout)
	}

	if strings.TrimSpace(result.Stdout) == "" {
		return []*GroupMember{}, nil
	}

	gm, err := unmarshalGroupMembership([]byte(result.Stdout))
	if err != nil {
		return nil, fmt.Errorf("while unmarshalling group membership response: %s", err)
	}

	return gm, nil
}

func (g *GroupMembership) bulkGroupMembersOp(conf *config.ProviderConf, operation string, members []*GroupMember) error {
	if len(members) == 0 {
		return nil
	}

	memberList := getMembershipList(members)
	cmd := fmt.Sprintf("%s -Identity %q %s -Confirm:$false", operation, g.GroupGUID, memberList)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      false,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)

	if err != nil {
		return fmt.Errorf("while running %s: %s", operation, err)
	} else if result.ExitCode != 0 {
		return fmt.Errorf("command %s exited with a non-zero exit code(%d), stderr: %s, stdout: %s", operation, result.ExitCode, result.StdErr, result.Stdout)
	}

	return nil
}

func (g *GroupMembership) AddGroupMembers(conf *config.ProviderConf, members []*GroupMember) error {
	return g.bulkGroupMembersOp(conf, "Add-ADGroupMember", members)
}

func (g *GroupMembership) RemoveGroupMembers(conf *config.ProviderConf, members []*GroupMember) error {
	return g.bulkGroupMembersOp(conf, "Remove-ADGroupMember", members)
}

func (g *GroupMembership) SetGroupMembers(conf *config.ProviderConf, expected []*GroupMember) error {
	existing, err := g.getGroupMembers(conf)
	if err != nil {
		return err
	}

	toAdd, toRemove := diffGroupMemberLists(expected, existing)
	err = g.AddGroupMembers(conf, toAdd)
	if err != nil {
		return err
	}

	err = g.RemoveGroupMembers(conf, toRemove)
	if err != nil {
		return err
	}

	return nil
}

func (g *GroupMembership) Create(conf *config.ProviderConf) error {
	if len(g.GroupMembers) == 0 {
		return nil
	}

	memberList := getMembershipList(g.GroupMembers)
	cmds := []string{fmt.Sprintf("Add-ADGroupMember -Identity %q -Members %s", g.GroupGUID, memberList)}
	psOpts := CreatePSCommandOpts{
		JSONOutput:      false,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
	}
	psCmd := NewPSCommand(cmds, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return fmt.Errorf("while running Add-ADGroupMember: %s", err)
	} else if result.ExitCode != 0 {
		return fmt.Errorf("command Add-ADGroupMember exited with a non-zero exit code(%d), stderr: %s, stdout: %s", result.ExitCode, result.StdErr, result.Stdout)
	}

	return nil
}

func (g *GroupMembership) Delete(conf *config.ProviderConf) error {
	subCmdOpt := CreatePSCommandOpts{
		JSONOutput:      false,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
		SkipCredPrefix:  true,
	}
	subcmd := NewPSCommand([]string{fmt.Sprintf("Get-AdGroupMember %q", g.GroupGUID)}, subCmdOpt)
	cmd := fmt.Sprintf("Remove-ADGroupMember %q -Members (%s) -Confirm:$false", g.GroupGUID, subcmd.String())

	psOpts := CreatePSCommandOpts{
		JSONOutput:      false,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
	if err != nil {
		return fmt.Errorf("while running Remove-ADGroupMember: %s", err)
	} else if result.ExitCode != 0 && !strings.Contains(result.StdErr, "InvalidData") {
		return fmt.Errorf("command Remove-ADGroupMember exited with a non-zero exit code(%d), stderr: %s, stdout: %s", result.ExitCode, result.StdErr, result.Stdout)
	}
	return nil
}

func NewGroupMembershipFromHost(conf *config.ProviderConf, groupID string) (*GroupMembership, error) {
	result := &GroupMembership{
		GroupGUID: groupID,
	}

	gm, err := result.getGroupMembers(conf)
	if err != nil {
		return nil, err
	}
	result.GroupMembers = gm

	return result, nil
}

func NewGroupMembershipFromState(d *schema.ResourceData) (*GroupMembership, error) {
	groupID := d.Get("group_id").(string)
	members := d.Get("group_members").(*schema.Set)
	result := &GroupMembership{
		GroupGUID:    groupID,
		GroupMembers: []*GroupMember{},
	}

	for _, m := range members.List() {
		if m == "" {
			continue
		}
		newMember := &GroupMember{
			GUID: m.(string),
		}

		result.GroupMembers = append(result.GroupMembers, newMember)
	}
	return result, nil
}
