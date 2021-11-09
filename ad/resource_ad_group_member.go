package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGroupMember() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_group_member` manages a specific member of a given Active Directory group.",
		Create:      resourceADGroupMemberCreate,
		Read:        resourceADGroupMemberRead,
		Delete:      resourceADGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the group. This can be a GUID, a SID, a Distinguished Name, or the SAM Account Name of the group.",
				ForceNew:    true,
			},
			"group_member": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A member AD Principal. The principal can be identified by its GUID, SID, Distinguished Name, or SAM Account Name.",
				ForceNew:    true,
			},
		},
	}
}

func composeGroupMemberID(groupID, memberID string) string {
	return groupID + "_" + memberID
}

func parseGroupMemberID(groupMemberID string) (groupID, memberID string, err error) {
	ids := strings.Split(groupMemberID, "_")

	if len(ids) != 2 {
		err = fmt.Errorf("invalid groupMemberID: %s", groupMemberID)
		return
	}

	groupID = ids[0]
	memberID = ids[1]

	return
}

func resourceADGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	groupID, memberID, err := parseGroupMemberID(d.Id())
	if err != nil {
		// This is a provider internal error. Let's return it.
		return err
	}

	gm, err := winrmhelper.NewGroupMembershipFromHost(meta.(*config.ProviderConf), groupID)
	if err != nil {
		return err
	}

	for _, m := range gm.GroupMembers {
		if memberID == m.GUID {

			_ = d.Set("group_member", memberID)
			_ = d.Set("group_id", groupID)

			return nil
		}
	}

	log.Printf("error finding member %s in membership of group %s", memberID, groupID)
	d.SetId("")
	return nil
}

func resourceADGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	groupID := d.Get("group_id").(string)
	memberID := d.Get("group_member").(string)

	gm := &winrmhelper.GroupMembership{
		GroupGUID: groupID,
		GroupMembers: []*winrmhelper.GroupMember{
			{
				GUID: memberID,
			},
		},
	}

	err := gm.Create(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}

	d.SetId(composeGroupMemberID(groupID, memberID))

	return nil
}

func resourceADGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	groupID, memberID, err := parseGroupMemberID(d.Id())
	if err != nil {
		// This is a provider internal error. Let's return it.
		return err
	}

	gm := &winrmhelper.GroupMembership{
		GroupGUID: groupID,
		GroupMembers: []*winrmhelper.GroupMember{
			{
				GUID: memberID,
			},
		},
	}

	err = gm.RemoveGroupMembers(meta.(*config.ProviderConf), gm.GroupMembers)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
