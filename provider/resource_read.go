package provider

import (
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/state"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRead(driverName string) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*libmachine.Client)
		name := d.Get("name").(string)
		var err error
		h, err := client.Load(name)
		if err != nil {
			return err
		}
		machineState, err := h.Driver.GetState()
		if err != nil {
			return fmt.Errorf("Error attempting to retrieve state: %s", err)
		}
		if machineState == state.Running {
			sshHostname, err := h.Driver.GetSSHHostname()
			if err != nil {
				return fmt.Errorf("Error attempting to retrieve ssh hostname: %s", err)
			}
			d.Set("ssh_hostname", sshHostname)
			sshPort, err := h.Driver.GetSSHPort()
			if err != nil {
				return fmt.Errorf("Error attempting to retrieve ssh port: %s", err)
			}
			d.Set("ssh_port", sshPort)
			address, err := h.Driver.GetIP()
			if err != nil {
				return fmt.Errorf("Error attempting to retrieve address: %s", err)
			}
			d.Set("address", address)
			dockerUrl, err := h.Driver.GetURL()
			if err != nil {
				return fmt.Errorf("Error attempting to retrieve docker url: %s", err)
			}
			d.Set("docker_url", dockerUrl)
			dockerVersion, err := h.DockerVersion()
			if err != nil {
				return fmt.Errorf("Error attempting to retrieve docker version: %s", err)
			}
			d.Set("docker_version", dockerVersion)
		} else {
			d.Set("ssh_hostname", nil)
			d.Set("ssh_port", nil)
			d.Set("address", nil)
			d.Set("docker_url", nil)
			d.Set("docker_version", nil)
		}
		d.Set("state", strings.ToLower(machineState.String()))
		return nil
	}
}
