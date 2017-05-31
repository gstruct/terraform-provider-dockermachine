package provider

import (
	"fmt"

	"github.com/docker/machine/libmachine"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDelete(driverName string) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*libmachine.Client)
		name := d.Get("name").(string)
		var err error
		host, err := client.Load(name)
		if err != nil {
			return err
		}
		err = host.Driver.Remove()
		if err != nil {
			return fmt.Errorf("Error removing host %q: %s", name, err)
		}

		exist, err := client.Exists(name)
		if err != nil {
			return fmt.Errorf("Error removing host %q: %s", name, err)
		}
		if !exist {
			return fmt.Errorf("Error removing host %q: host does not exist.", name)
		}
		return client.Remove(name)
	}
}
