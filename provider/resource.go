package provider

import (
	"strings"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/drivers/rpc"
	"github.com/docker/machine/libmachine/mcnflag"

	"github.com/docker/machine/drivers/amazonec2"
	"github.com/docker/machine/drivers/azure"
	"github.com/docker/machine/drivers/digitalocean"
	"github.com/docker/machine/drivers/exoscale"
	"github.com/docker/machine/drivers/generic"
	"github.com/docker/machine/drivers/google"
	"github.com/docker/machine/drivers/hyperv"
	"github.com/docker/machine/drivers/none"
	"github.com/docker/machine/drivers/openstack"
	"github.com/docker/machine/drivers/rackspace"
	"github.com/docker/machine/drivers/softlayer"
	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/drivers/vmwarefusion"
	"github.com/docker/machine/drivers/vmwarevcloudair"
	"github.com/docker/machine/drivers/vmwarevsphere"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resource(driverName string) *schema.Resource {
	drv := getDriver(driverName, "", "")
	resourceSchema := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"certs_directory": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"tls_ca_cert": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tls_ca_key": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tls_client_cert": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tls_client_key": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tls_server_cert": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tls_server_key": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"storage_path": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"storage_path_computed": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"tls_san": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_opt": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_env": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_insecure_registry": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_label": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_registry_mirror": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"engine_storage_driver": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"engine_install_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		"swarm_master": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		"swarm_image": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm_discovery": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm_addr": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm_host": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm_strategy": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"swarm_opt": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"swarm_join_opt": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			ForceNew: true,
		},
		"swarm_experimental": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},
		"ssh_hostname": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssh_port": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ssh_username": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssh_keypath": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"address": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"docker_version": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "running",
			ValidateFunc: validation.StringInSlice([]string{"running", "stopped"}, false),
		},
	}
	for _, flag := range drv.GetCreateFlags() {
		flagName := strings.Replace(flag.String(), "-", "_", -1)
		switch f := flag.(type) {
		case mcnflag.StringFlag:
			resourceSchema[flagName] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  f.Value,
			}
		case mcnflag.StringSliceFlag:
			resourceSchema[flagName] = &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			}
		case mcnflag.IntFlag:
			resourceSchema[flagName] = &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  f.Value,
			}
		case mcnflag.BoolFlag:
			resourceSchema[flagName] = &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			}
		}
	}
	return &schema.Resource{
		Schema: resourceSchema,
		Exists: resourceExists(drv.DriverName()),
		Create: resourceCreate(drv.DriverName()),
		Read:   resourceRead(drv.DriverName()),
		Update: resourceUpdate(drv.DriverName()),
		Delete: resourceDelete(drv.DriverName()),
	}
}

func getDriver(driverName, machineName, storePath string) drivers.Driver {
	switch driverName {
	case "amazonec2":
		return amazonec2.NewDriver(machineName, storePath)
	case "azure":
		return azure.NewDriver(machineName, storePath)
	case "digitalocean":
		return digitalocean.NewDriver(machineName, storePath)
	case "exoscale":
		return exoscale.NewDriver(machineName, storePath)
	case "generic":
		return generic.NewDriver(machineName, storePath)
	case "google":
		return google.NewDriver(machineName, storePath)
	case "hyperv":
		return hyperv.NewDriver(machineName, storePath)
	case "none":
		return none.NewDriver(machineName, storePath)
	case "openstack":
		return openstack.NewDriver(machineName, storePath)
	case "rackspace":
		return rackspace.NewDriver(machineName, storePath)
	case "softlayer":
		return softlayer.NewDriver(machineName, storePath)
	case "virtualbox":
		return virtualbox.NewDriver(machineName, storePath)
	case "vmwarefusion":
		return vmwarefusion.NewDriver(machineName, storePath)
	case "vmwarevcloudair":
		return vmwarevcloudair.NewDriver(machineName, storePath)
	case "vmwarevsphere":
		return vmwarevsphere.NewDriver(machineName, storePath)
	default:
		return nil
	}
}

func getDriverOpts(d *schema.ResourceData, mcnflags []mcnflag.Flag) drivers.DriverOptions {
	driverOpts := rpcdriver.RPCFlags{
		Values: make(map[string]interface{}),
	}

	for _, f := range mcnflags {
		driverOpts.Values[f.String()] = f.Default()

		if f.Default() == nil {
			driverOpts.Values[f.String()] = false
		}

		schemaOpt := strings.Replace(f.String(), "-", "_", -1)
		switch f.(type) {
		case mcnflag.StringFlag:
			driverOpts.Values[f.String()] = d.Get(schemaOpt).(string)
		case mcnflag.StringSliceFlag:
			driverOpts.Values[f.String()] = ss2is(d.Get(schemaOpt).([]string))
		case mcnflag.IntFlag:
			driverOpts.Values[f.String()] = d.Get(schemaOpt).(int)
		case mcnflag.BoolFlag:
			driverOpts.Values[f.String()] = d.Get(schemaOpt).(bool)
		}
	}

	return driverOpts
}
