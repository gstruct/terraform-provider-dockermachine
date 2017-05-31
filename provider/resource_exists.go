package provider

import (
	"github.com/docker/machine/libmachine"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceExists(driverName string) func(*schema.ResourceData, interface{}) (bool, error) {
	return func(d *schema.ResourceData, meta interface{}) (bool, error) {
		client := meta.(*libmachine.Client)
		name := d.Get("name").(string)
		exists, err := client.Exists(name)
		if err != nil {
			return false, err
		}
		if !exists {
			d.SetId(name)
			return false, nil
		}
		return true, nil
	}
}
