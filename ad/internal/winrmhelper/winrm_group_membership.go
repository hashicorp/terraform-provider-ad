package winrmhelper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
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

	return gm, nil
}

func getMembershipList(g []*GroupMember) string {
	out := []string{}
	for _, member := range g {
		out = append(out, member.GUID)
	}

	return strings.Join(out, ",")
}

func (g *GroupMembership) getGroupMembers(client *winrm.Client) ([]*GroupMember, error) {
	cmd := fmt.Sprintf("Get-ADGroupMember -Identity %q", g.GroupGUID)

	result, err := RunWinRMCommand(client, []string{cmd}, true, true)
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

func (g *GroupMembership) bulkGroupMembersOp(client *winrm.Client, operation string, members []*GroupMember) error {
	if len(members) == 0 {
		return nil
	}

	memberList := getMembershipList(members)
	cmd := fmt.Sprintf("%s -Identity %q %s -Confirm:$false", operation, g.GroupGUID, memberList)

	result, err := RunWinRMCommand(client, []string{cmd}, false, false)
	if err != nil {
		return fmt.Errorf("while running %s: %s", operation, err)
	} else if result.ExitCode != 0 {
		return fmt.Errorf("command %s exited with a non-zero exit code(%d), stderr: %s, stdout: %s", operation, result.ExitCode, result.StdErr, result.Stdout)
	}

	return nil
}

func (g *GroupMembership) addGroupMembers(client *winrm.Client, members []*GroupMember) error {
	return g.bulkGroupMembersOp(client, "Add-ADGroupMember", members)
}

func (g *GroupMembership) removeGroupMembers(client *winrm.Client, members []*GroupMember) error {
	return g.bulkGroupMembersOp(client, "Remove-ADGroupMember", members)
}

func (g *GroupMembership) Update(client *winrm.Client, expected []*GroupMember) error {
	existing, err := g.getGroupMembers(client)
	if err != nil {
		return err
	}

	toAdd, toRemove := diffGroupMemberLists(expected, existing)
	err = g.addGroupMembers(client, toAdd)
	if err != nil {
		return err
	}

	err = g.removeGroupMembers(client, toRemove)
	if err != nil {
		return err
	}

	return nil
}

func (g *GroupMembership) Create(client *winrm.Client) error {
	if len(g.GroupMembers) == 0 {
		return nil
	}

	memberList := getMembershipList(g.GroupMembers)
	cmd := []string{fmt.Sprintf("Add-ADGroupMember -Identity %q -Members %s", g.GroupGUID, memberList)}
	result, err := RunWinRMCommand(client, cmd, false, false)
	if err != nil {
		return fmt.Errorf("while running Add-ADGroupMember: %s", err)
	} else if result.ExitCode != 0 {
		return fmt.Errorf("command Add-ADGroupMember exited with a non-zero exit code(%d), stderr: %s, stdout: %s", result.ExitCode, result.StdErr, result.Stdout)
	}

	return nil
}

func (g *GroupMembership) Delete(client *winrm.Client) error {
	cmd := fmt.Sprintf("Remove-ADGroupMember %q -Members (Get-ADGroupMember %q) -Confirm:$false", g.GroupGUID, g.GroupGUID)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false)
	if err != nil {
		return fmt.Errorf("while running Remove-ADGroupMember: %s", err)
	} else if result.ExitCode != 0 && !strings.Contains(result.StdErr, "InvalidData") {
		return fmt.Errorf("command Remove-ADGroupMember exited with a non-zero exit code(%d), stderr: %s, stdout: %s", result.ExitCode, result.StdErr, result.Stdout)
	}
	return nil
}

func NewGroupMembershipFromHost(client *winrm.Client, groupID string) (*GroupMembership, error) {
	result := &GroupMembership{
		GroupGUID: groupID,
	}

	gm, err := result.getGroupMembers(client)
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
