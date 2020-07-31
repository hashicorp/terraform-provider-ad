package ad

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/adschema"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/gposec"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/winrmhelper"
)

func resourceADGPOSecurity() *schema.Resource {
	return &schema.Resource{
		Description: "`ad_gpo_security` manages the security settings portion of a Group Policy Object (GPO).",
		Create:      resourceADGPOSecurityCreate,
		Read:        resourceADGPOSecurityRead,
		Update:      resourceADGPOSecurityUpdate,
		Delete:      resourceADGPOSecurityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: adschema.GpoSecuritySchema(),
	}
}

func resourceADGPOSecurityCreate(d *schema.ResourceData, meta interface{}) error {
	winrmClient := meta.(ProviderConf).WinRMClient
	winrmCPClient := meta.(ProviderConf).WinRMCPClient

	guid := d.Get("gpo_container").(string)
	if guid == "" {
		return fmt.Errorf("Cannot handle empty GPO GUID")
	}
	_, err := uuid.ParseUUID(guid)
	if err != nil {
		return fmt.Errorf("Cannot parse GUID %q: %s", guid, err)
	}
	iniFile, err := winrmhelper.GetSecIniFromResource(d, adschema.GpoSecuritySchema())
	if err != nil {
		return fmt.Errorf("error while generating ini file from resource data: %s", err)
	}

	gpo, err := winrmhelper.GetGPOFromHost(winrmClient, "", guid)
	if err != nil {
		return err
	}

	err = winrmhelper.UploadSecIni(winrmClient, winrmCPClient, gpo, iniFile)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_securitysettings", guid))
	return resourceADGPOSecurityRead(d, meta)
}

func resourceADGPOSecurityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ProviderConf).WinRMClient
	resourceID := d.Id()
	toks := strings.Split(resourceID, "_")
	if len(toks) != 2 {
		return fmt.Errorf("resource ID %q does not match <guid>_securitysettings", resourceID)
	}
	guid := toks[0]

	gpo, err := winrmhelper.GetGPOFromHost(client, "", guid)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			log.Printf("[DEBUG] GPO with guid %q not found", guid)
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("gpo_container", guid)

	hostSecIni, err := winrmhelper.GetSecIniFromHost(client, gpo)
	if err != nil {
		if strings.Contains(err.Error(), "ItemNotFoundException") {
			log.Printf("[DEBUG] inf file not found, marking resource as gone")
			d.SetId("")
			return nil
		}
		return err
	}

	err = gposec.HandleSectionRead(adschema.GPOSecuritySchemaKeys, hostSecIni, d)
	return err
}

func resourceADGPOSecurityUpdate(d *schema.ResourceData, meta interface{}) error {
	winrmClient := meta.(ProviderConf).WinRMClient
	winrmCPClient := meta.(ProviderConf).WinRMCPClient

	guid := d.Get("gpo_container").(string)
	if guid == "" {
		return fmt.Errorf("Cannot handle empty GPO GUID")
	}
	_, err := uuid.ParseUUID(guid)
	if err != nil {
		return fmt.Errorf("Cannot parse GUID %q: %s", guid, err)
	}

	gpo, err := winrmhelper.GetGPOFromHost(winrmClient, "", guid)
	if err != nil {
		return fmt.Errorf("error while retrieving GPO with guid %q: %s", guid, err)
	}

	iniFile, err := winrmhelper.GetSecIniFromResource(d, adschema.GpoSecuritySchema())
	if err != nil {
		return fmt.Errorf("error while generating ini file from resource data: %s", err)
	}

	iniBuf := bytes.NewBuffer([]byte{})
	_, err = iniFile.WriteTo(iniBuf)
	if err != nil {
		return fmt.Errorf("error while writing INI file in buffer")
	}
	iniSum := sha256.Sum256(iniBuf.Bytes())

	hostSecIniBytes, err := winrmhelper.GetSecIniContents(winrmClient, gpo)
	if err != nil {
		return fmt.Errorf("error while retrieving security settings contents for GPO with guid %q: %s", guid, err)
	}

	hostSum := sha256.Sum256(hostSecIniBytes)

	if iniSum != hostSum {
		err = winrmhelper.UploadSecIni(winrmClient, winrmCPClient, gpo, iniFile)
		if err != nil {
			return fmt.Errorf("error while uploading security settings file for GPO with guid %q: %s", guid, err)
		}

	}
	return resourceADGPOSecurityRead(d, meta)
}

func resourceADGPOSecurityDelete(d *schema.ResourceData, meta interface{}) error {
	winrmClient := meta.(ProviderConf).WinRMClient
	winrmCPClient := meta.(ProviderConf).WinRMCPClient
	resourceID := d.Id()
	toks := strings.Split(resourceID, "_")
	if len(toks) != 2 {
		return fmt.Errorf("resource ID %q does not match <guid>_securitysettings", resourceID)
	}
	guid := toks[0]

	gpo, err := winrmhelper.GetGPOFromHost(winrmClient, "", guid)
	if err != nil {
		return fmt.Errorf("error while retrieving GPO with guid %q: %s", guid, err)
	}

	err = winrmhelper.RemoveSecIni(winrmClient, winrmCPClient, gpo)
	if err != nil {
		return fmt.Errorf("error while removing security settings INF file for GPO with guid %q: %s", guid, err)
	}
	return resourceADGPOSecurityRead(d, meta)
}
