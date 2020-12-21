package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADUser() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_user` manages User objects in an Active Directory tree.",
		Create:      resourceADUserCreate,
		Read:        resourceADUserRead,
		Update:      resourceADUserUpdate,
		Delete:      resourceADUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Display Name of an Active Directory user.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Principal Name of an Active Directory user.",
			},
			"sam_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pre-win2k user logon name.",
			},
			"initial_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user's initial password. This will be set on creation but will *not* be enforced in subsequent plans.",
			},
			"container": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "A DN of the container object that will be holding the user.",
				DiffSuppressFunc: suppressCaseDiff,
			},
			"cannot_change_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the user will not be allowed to change their password.",
			},
			"password_never_expires": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the password for this user will not expire.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to false, the user will be disabled.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's town or city. This parameter sets the City property of a user object.",
			},
			"company": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's company. This parameter sets the Company property of a user object.",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the country or region code for the user's language of choice. This parameter sets the Country property of a user object.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's department. This parameter sets the Department property of a user object.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description of the object. This parameter sets the value of the Description property for the user object.",
			},
			"division": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's division. This parameter sets the Division property of a user object.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's e-mail address. This parameter sets the EmailAddress property of a user object.",
			},
			"employee_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's employee ID. This parameter sets the EmployeeID property of a user object.",
			},
			"employee_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's employee number. This parameter sets the EmployeeNumber property of a user object.",
			},
			"fax": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's fax phone number. This parameter sets the Fax property of a user object.",
			},
			"given_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's given name. This parameter sets the GivenName property of a user object.",
			},
			"home_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a user's home directory. This parameter sets the HomeDirectory property of a user object.",
			},
			"home_drive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a drive that is associated with the UNC path defined by the HomeDirectory property. The drive letter is specified as <DriveLetter>: where <DriveLetter> indicates the letter of the drive to associate. The <DriveLetter> must be a single, uppercase letter and the colon is required. This parameter sets the HomeDrive property of the user object.",
			},
			"home_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's home telephone number. This parameter sets the HomePhone property of a user object.",
			},
			"home_page": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the URL of the home page of the object. This parameter sets the homePage property of a user object.",
			},
			"initials": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the initials that represent part of a user's name. You can use this value for the user's middle initial. This parameter sets the Initials property of a user object.",
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's mobile phone number. This parameter sets the MobilePhone property of a user object.",
			},
			"office": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the location of the user's office or place of business. This parameter sets the Office property of a user object.",
			},
			"office_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's office telephone number. This parameter sets the OfficePhone property of a user object.",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's organization. This parameter sets the Organization property of a user object.",
			},
			"other_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a name in addition to a user's given name and surname, such as the user's middle name.",
			},
			"po_box": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's post office box number. This parameter sets the POBox property of a user object.",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's postal code or zip code. This parameter sets the PostalCode property of a user object.",
			},
			"smart_card_logon_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, a smart card is required to logon. This parameter sets the SmartCardLoginRequired property for a user object.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's or Organizational Unit's state or province. This parameter sets the State property of a user object.",
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's street address. This parameter sets the StreetAddress property of a user object.",
			},
			"surname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's last name or surname. This parameter sets the Surname property of a user object.",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's title. This parameter sets the Title property of a user object",
			},
			"trusted_for_delegation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the user account is trusted for Kerberos delegation. A service that runs under an account that is trusted for Kerberos delegation can assume the identity of a client requesting the service. This parameter sets the TrustedForDelegation property of an account object.",
			},
		},
	}
}

func resourceADUserCreate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetUserFromResource(d)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	guid, err := u.NewUser(client)
	if err != nil {
		return err
	}
	d.SetId(guid)
	return resourceADUserRead(d, meta)
}

func resourceADUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading ad_user resource for user with guid: %q", d.Id())
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			d.SetId("")
			return nil
		}
		return err
	}
	if u == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("sam_account_name", u.SAMAccountName)
	_ = d.Set("display_name", u.DisplayName)
	_ = d.Set("principal_name", u.PrincipalName)
	_ = d.Set("container", u.Container)
	_ = d.Set("enabled", u.Enabled)
	_ = d.Set("password_never_expires", u.PasswordNeverExpires)
	_ = d.Set("cannot_change_password", u.CannotChangePassword)
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

	return nil
}

func resourceADUserUpdate(d *schema.ResourceData, meta interface{}) error {
	u := winrmhelper.GetUserFromResource(d)
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	err = u.ModifyUser(d, client)
	if err != nil {
		return err
	}
	return resourceADUserRead(d, meta)
}

func resourceADUserDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(ProviderConf).AcquireWinRMClient()
	if err != nil {
		return err
	}
	defer meta.(ProviderConf).ReleaseWinRMClient(client)

	u, err := winrmhelper.GetUserFromHost(client, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "ADIdentityNotFoundException") {
			return nil
		}
		return fmt.Errorf("while retrieving user data from host: %s", err)
	}
	err = u.DeleteUser(client)
	if err != nil {
		return fmt.Errorf("while deleting user: %s", err)
	}
	return resourceADUserRead(d, meta)
}
