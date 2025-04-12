package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGroupMembership() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_group_membership` manages the members of a given Active Directory group.",
		Create:      resourceADGroupMembershipCreate,
		Read:        resourceADGroupMembershipRead,
		Update:      resourceADGroupMembershipUpdate,
		Delete:      resourceADGroupMembershipDelete,
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
			"group_members": {
				Type:             schema.TypeSet,
				Required:         true,
				Description:      "A list of member AD Principals. Each principal can be identified by its GUID, SID, Distinguished Name, or SAM Account Name. Only one is required",
				Elem:             &schema.Schema{Type: schema.TypeString},
				MinItems:         1,
				DiffSuppressFunc: suppressGroupMemberDiff,
			},
			"group_members_details": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sam_account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Detailed attributes of group members, resolved internally by the provider.",
			},
		},
	}
}

func resourceADGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	toks := strings.Split(d.Id(), "_")

	gm, err := winrmhelper.NewGroupMembershipFromHost(meta.(*config.ProviderConf), toks[0])
	if err != nil {
		return err
	}
	memberList := d.Get("group_members").(*schema.Set).List()
	memberDetails := []map[string]string{}

	for _, identity := range memberList {
		id := strings.Trim(identity.(string), "\"")
		for _, mbr := range gm.GroupMembers {
			if id == mbr.GUID || id == mbr.DN || id == mbr.SamAccountName || id == mbr.SID.Value {
				// Resolve additional details for each member
				details := map[string]string{
					"guid":             mbr.GUID,
					"dn":               mbr.DN,
					"sam_account_name": mbr.SamAccountName,
					"sid":              mbr.SID.Value,
					"identity":         identity.(string),
				}
				memberDetails = append(memberDetails, details)
				break
			}
		}
	}

	_ = d.Set("group_members", memberList)
	_ = d.Set("group_members_details", memberDetails)
	_ = d.Set("group_id", toks[0])
	return nil
}

func resourceADGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	gm, err := winrmhelper.NewGroupMembershipFromState(d)
	if err != nil {
		return err
	}

	err = gm.Create(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}

	membershipUUID, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("while generating UUID to use as unique membership ID: %s", err)
	}

	id := fmt.Sprintf("%s_%s", gm.GroupGUID, membershipUUID)
	d.SetId(id)

	return nil
}

func resourceADGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	gm, err := winrmhelper.NewGroupMembershipFromState(d)
	if err != nil {
		return err
	}

	err = gm.Update(meta.(*config.ProviderConf), gm.GroupMembers)
	if err != nil {
		return err
	}

	return resourceADGroupMembershipRead(d, meta)
}

func resourceADGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	gm, err := winrmhelper.NewGroupMembershipFromState(d)
	if err != nil {
		return err
	}

	err = gm.Delete(meta.(*config.ProviderConf))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func suppressGroupMemberDiff(k, old, new string, d *schema.ResourceData) bool {
	strippedNew := strings.Trim(new, "\"")

	// Get the resolved member details from the state
	memberDetails := d.Get("group_members_details").(*schema.Set).List()

	for _, member := range memberDetails {
		// Check if the new value matches any of the saved group member details
		for _, member := range member.(map[string]string) {
			if member == strippedNew {
				// Match found, suppress the diff
				return true
			}
		}
	}

	return false
}
