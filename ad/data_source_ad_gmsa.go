package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceADGmsa() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory gmsa object.",
		Read:        dataSourceADgmsaRead,
		Schema: map[string]*schema.Schema{
			"gmsa_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The gmsa's identifier. It can be the group's GUID, SID, Distinguished Name, or SAM Account Name.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Display Name of an Active Directory Gmsa.",
			},
			"delegated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to false, the Gmsa will not be delegated to a service.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies a description of the object. This parameter sets the value of the Description property for the Gmsa object.",
			},
			"dns_host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the dns host name of the Gmsa object.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to false, the Gmsa will be disabled.",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration date of the gmsa.",
			},
			"home_page": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the URL of the home page of the object. This parameter sets the homePage property of a Gmsa object.",
			},
			"managed_password_interval_in_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the number of days for the password change interval.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the name of the Gmsa object.",
			},
			"principals_allowed_to_delegate_to_account": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "This value sets the encryption types supported flags of the Active Directory.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"principals_allowed_to_retrieve_managed_password": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "This value sets the encryption types supported flags of the Active Directory.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The pre-win2k Gmsa logon name.",
			},
			"service_principal_names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "This value sets SPN's for the gmsa.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"trusted_for_delegation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, the Gmsa account is trusted for Kerberos delegation. A service that runs under an account that is trusted for Kerberos delegation can assume the identity of a client requesting the service. This parameter sets the TrustedForDelegation property of an account object.",
			},
			"sid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SID of the gmsa object.",
			},
		},
	}
}

func dataSourceADgmsaRead(d *schema.ResourceData, meta interface{}) error {
	isLocal := meta.(ProviderConf).isConnectionTypeLocal()
	gmsaID := d.Get("gmsa_id").(string)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	g, err := winrmhelper.GetGmsaFromHost(client, gmsaID, isLocal)
	if err != nil {
		return err
	}

	if g == nil {
		return fmt.Errorf("No gmsa found with gmsa_id %q", gmsaID)
	}
	_ = d.Set("container", g.Container)
	_ = d.Set("display_name", g.DisplayName)
	_ = d.Set("delegated", g.Delegated)
	_ = d.Set("description", g.Description)
	_ = d.Set("dns_host_name", g.DNSHostName)
	_ = d.Set("enabled", g.Enabled)
	_ = d.Set("expiration", g.Expiration)
	_ = d.Set("home_page", g.HomePage)
	_ = d.Set("KerberosEncryptionType", g.KerberosEncryptionType)
	_ = d.Set("managed_password_interval_in_days", g.ManagedPasswordIntervalInDays)
	_ = d.Set("name", g.Name)
	_ = d.Set("principals_allowed_to_delegate_to_account", g.PrincipalsAllowedToDelegateToAccount)
	_ = d.Set("principals_allowed_to_retrieve_managed_password", g.PrincipalsAllowedToRetrieveManagedPassword)
	_ = d.Set("sam_account_name", g.SAMAccountName)
	_ = d.Set("service_principal_names", g.ServicePrincipalNames)
	_ = d.Set("trusted_for_delegation", g.TrustedForDelegation)
	_ = d.Set("sid", g.SID.Value)
	d.SetId(g.GUID)

	return nil
}
