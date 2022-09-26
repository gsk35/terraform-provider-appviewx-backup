package appviewx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-appviewx/appviewx/config"
	"terraform-provider-appviewx/appviewx/constants"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			constants.APPVIEWX_USERNAME: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_PASSWORD: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_IP: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_PORT: {
				Type:     schema.TypeString,
				Required: true,
			},
			constants.APPVIEWX_ENVIRONMENT_Is_HTTPS: {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appviewx_automation":  ResourceAutomationServer(),
			"appviewx_certificate": ResourceCertificateServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	appviewxEnvironment := config.AppViewXEnvironment{
		AppViewXUserName:        d.Get(constants.APPVIEWX_USERNAME).(string),
		AppViewXPassword:        d.Get(constants.APPVIEWX_PASSWORD).(string),
		AppViewXEnvironmentIP:   d.Get(constants.APPVIEWX_ENVIRONMENT_IP).(string),
		AppViewXEnvironmentPort: d.Get(constants.APPVIEWX_ENVIRONMENT_PORT).(string),
		AppViewXIsHTTPS:         d.Get(constants.APPVIEWX_ENVIRONMENT_Is_HTTPS).(bool),
	}
	return &appviewxEnvironment, nil
}
