package msad

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Config holds all the information necessary to configure the provider
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Protocol string
	Insecure bool
}

// NewConfig returns a new Config struct populated with Resource Data.
func NewConfig(d *schema.ResourceData) (*Config, error) {
	return nil, nil
}
