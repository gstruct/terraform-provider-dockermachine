package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/drivers/plugin/localbinary"
	"github.com/docker/machine/libmachine/log"
)

func Provider() terraform.ResourceProvider {
	resourceMap := make(map[string]*schema.Resource)
	for _, str := range localbinary.CoreDrivers {
		resourceMap[fmt.Sprintf("dockermachine_%s", str)] = resource(str)
	}
	//resourceMap["dockermachine_external"] = resourceExternal()
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		ResourcesMap:  resourceMap,
		Schema: map[string]*schema.Schema{
			"debug": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "docker-machine debug output",
			},
			"storage_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: storagePathDefault,
				Description: "docker-machine storage path",
			},
			"certs_directory": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: certsDirDefault,
				Description: "docker-machine certificates directory",
			},
		},
	}
}

func storagePathDefault() (interface{}, error) {
	return mcndirs.GetBaseDir(), nil
}

func certsDirDefault() (interface{}, error) {
	return mcndirs.GetMachineCertDir(), nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.SetDebug(d.Get("debug").(bool))
	return libmachine.NewClient(d.Get("storage_path").(string), d.Get("certs_directory").(string)), nil
}
