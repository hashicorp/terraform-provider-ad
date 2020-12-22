package ad

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceADUser() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an Active Directory user object.",
		Read:        dataSourceADUserRead,
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the user object.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAM account name of the user object.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the user object.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The principal name of the user object.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City assigned to user object.",
			},
			"company": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Company assigned to user object.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country assigned to user object.",
			},
			"department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Department assigned to user object.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the user object.",
			},
			"division": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Division assigned to user object.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address assigned to user object.",
			},
			"employee_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Employee ID assigned to user object.",
			},
			"employee_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Employee Number assigned to user object.",
			},
			"fax": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fax number assigned to user object.",
			},
			"given_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Given name of the user object.",
			},
			"home_directory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home directory of the user object.",
			},
			"home_drive": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home drive of the user object.",
			},
			"home_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home phone of the user object.",
			},
			"home_page": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Home page of the user object.",
			},
			"initials": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Initials of the user object.",
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mobile phone of the user object.",
			},
			"office": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Office assigned to user object.",
			},
			"office_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Office phone of the user object.",
			},
			"organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization assigned to user object.",
			},
			"other_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extra name of the user object.",
			},
			"po_box": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Post office assigned to user object.",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Postal code of the user object.",
			},
			"smart_card_logon_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Smart card required to logon or not",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the user object.",
			},
			"street_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address of the user object.",
			},
			"surname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Surname of the user object.",
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Title of the user object",
			},
			"trusted_for_delegation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Check if user is trusted for delegation",
			},
		},
	}
}

func dataSourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	dn := d.Get("guid").(string)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, dn)
	if err != nil {
		return err
	}

	if u == nil {
		return fmt.Errorf("No user found with dn %q", dn)
	}
	_ = d.Set("sam_account_name", u.SAMAccountName)
	_ = d.Set("display_name", u.DisplayName)
	_ = d.Set("principal_name", u.PrincipalName)
	_ = d.Set("guid", u.GUID)
	_ = d.Set("city", u.City)
	_ = d.Set("company", u.Company)
	_ = d.Set("country", u.Country)
	_ = d.Set("department", u.Department)
	_ = d.Set("description", u.Description)
	_ = d.Set("division", u.Division)
	_ = d.Set("email_address", u.EmailAddress)
	_ = d.Set("employee_id", u.EmployeeID)
	_ = d.Set("employee_number", u.EmployeeNumber)
	_ = d.Set("fax", u.Fax)
	_ = d.Set("given_name", u.GivenName)
	_ = d.Set("home_directory", u.HomeDirectory)
	_ = d.Set("home_drive", u.HomeDrive)
	_ = d.Set("home_phone", u.HomePhone)
	_ = d.Set("home_page", u.HomePage)
	_ = d.Set("initials", u.Initials)
	_ = d.Set("mobile_phone", u.MobilePhone)
	_ = d.Set("office", u.Office)
	_ = d.Set("office_phone", u.OfficePhone)
	_ = d.Set("organization", u.Organization)
	_ = d.Set("other_name", u.OtherName)
	_ = d.Set("po_box", u.POBox)
	_ = d.Set("postal_code", u.PostalCode)
	_ = d.Set("state", u.State)
	_ = d.Set("street_address", u.StreetAddress)
	_ = d.Set("surname", u.Surname)
	_ = d.Set("title", u.Title)
	_ = d.Set("smart_card_logon_required", u.SmartcardLogonRequired)
	_ = d.Set("trusted_for_delegation", u.TrustedForDelegation)
	d.SetId(u.GUID)

	return nil
}