package ad

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
	"github.com/packer-community/winrmcp/winrmcp"
)

// ProviderConfig holds all the information necessary to configure the provider
type ProviderConfig struct {
	WinRMUsername string
	WinRMPassword string
	WinRMHost     string
	WinRMPort     int
	WinRMProto    string
	WinRMInsecure bool
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) ProviderConfig {
	// winRM
	winRMUsername := d.Get("winrm_username").(string)
	winRMPassword := d.Get("winrm_password").(string)
	winRMHost := d.Get("winrm_hostname").(string)
	winRMPort := d.Get("winrm_port").(int)
	winRMProto := d.Get("winrm_proto").(string)
	winRMInsecure := d.Get("winrm_insecure").(bool)

	cfg := ProviderConfig{
		WinRMHost:     winRMHost,
		WinRMPort:     winRMPort,
		WinRMProto:    winRMProto,
		WinRMUsername: winRMUsername,
		WinRMPassword: winRMPassword,
		WinRMInsecure: winRMInsecure,
	}

	return cfg
}

// GetWinRMConnection returns a WinRM connection
func GetWinRMConnection(config ProviderConfig) (*winrm.Client, error) {
	useHTTPS := false
	if strings.ToLower(config.WinRMProto) == "https" {
		useHTTPS = true
	}

	endpoint := winrm.NewEndpoint(config.WinRMHost, config.WinRMPort, useHTTPS,
		config.WinRMInsecure, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, config.WinRMUsername, config.WinRMPassword)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetWinRMCPConnection sets up a winrmcp client that can be used to upload files to the DC.
func GetWinRMCPConnection(config ProviderConfig) (*winrmcp.Winrmcp, error) {
	useHTTPS := false
	if config.WinRMProto == "https" {
		useHTTPS = true
	}
	addr := fmt.Sprintf("%s:%d", config.WinRMHost, config.WinRMPort)
	cfg := winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     config.WinRMUsername,
			Password: config.WinRMPassword,
		},
		Https:                 useHTTPS,
		Insecure:              config.WinRMInsecure,
		MaxOperationsPerShell: 15,
	}
	return winrmcp.New(addr, &cfg)
}
