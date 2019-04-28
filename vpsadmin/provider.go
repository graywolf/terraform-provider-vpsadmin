package vpsadmin

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VPSADMIN_API_TOKEN", nil),
				Description: "The authentication token for API operations.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VPSADMIN_API_URL", "https://api.vpsfree.cz"),
				Description: "The URL to use for the vpsAdmin API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vpsadmin_ssh_key": resourceSshKey(),
			"vpsadmin_vps": resourceVps(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfg, err := configureClient(
		d.Get("api_url").(string),
		d.Get("auth_token").(string),
	)

	if err != nil {
		return nil, err
	}

	if err := cfg.testAuthentication(); err != nil {
		return nil, fmt.Errorf("Authentication failed: %v", err)
	}

	return cfg, nil
}
