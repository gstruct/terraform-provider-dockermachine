package provider

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/crashreport"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/docker/machine/libmachine/state"
	"github.com/docker/machine/libmachine/swarm"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCreate(driverName string) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*libmachine.Client)
		name := d.Get("name").(string)
		if !host.ValidateHostName(name) {
			return fmt.Errorf("Error creating machine: %s", mcnerror.ErrInvalidHostname)
		}
		drv := getDriver(driverName, name, client.Path)
		data, err := json.Marshal(drv)

		h, err := client.NewHost(drv.DriverName(), data)
		if err != nil {
			return err
		}

		storagePath := d.Get("storage_path").(string)
		if storagePath == "" {
			storagePath = filepath.Join(client.Path, "machines", name)
		}
		d.Set("storage_path_computed", storagePath)
		d.Set("tls_server_key", filepath.Join(storagePath, "server-key.pem"))
		d.Set("tls_server_cert", filepath.Join(storagePath, "server.pem"))

		certsDirectory := d.Get("certs_directory").(string)
		if certsDirectory == "" {
			certsDirectory = filepath.Join(client.Path, "certs")
		}

		h.HostOptions = &host.Options{
			AuthOptions: &auth.Options{
				CertDir:          certsDirectory,
				StorePath:        storagePath,
				ServerCertPath:   d.Get("tls_server_cert").(string),
				ServerKeyPath:    d.Get("tls_server_key").(string),
				CaCertPath:       tlsPath(d, "tls_ca_cert", certsDirectory, "ca.pem"),
				CaPrivateKeyPath: tlsPath(d, "tls_ca_key", certsDirectory, "ca-key.pem"),
				ClientCertPath:   tlsPath(d, "tls_client_cert", certsDirectory, "cert.pem"),
				ClientKeyPath:    tlsPath(d, "tls_client_key", certsDirectory, "key.pem"),
				ServerCertSANs:   is2ss(d.Get("tls_san").([]interface{})),
			},
			EngineOptions: &engine.Options{
				ArbitraryFlags:   is2ss(d.Get("engine_opt").([]interface{})),
				Env:              is2ss(d.Get("engine_env").([]interface{})),
				InsecureRegistry: is2ss(d.Get("engine_insecure_registry").([]interface{})),
				Labels:           is2ss(d.Get("engine_label").([]interface{})),
				RegistryMirror:   is2ss(d.Get("engine_registry_mirror").([]interface{})),
				StorageDriver:    d.Get("engine_storage_driver").(string),
				TLSVerify:        true,
				InstallURL:       d.Get("engine_install_url").(string),
			},
			SwarmOptions: &swarm.Options{
				IsSwarm:            d.Get("swarm").(bool) || d.Get("swarm_master").(bool),
				Image:              d.Get("swarm_image").(string),
				Agent:              d.Get("swarm").(bool),
				Master:             d.Get("swarm_master").(bool),
				Discovery:          d.Get("swarm_discovery").(string),
				Address:            d.Get("swarm_addr").(string),
				Host:               d.Get("swarm_host").(string),
				Strategy:           d.Get("swarm_strategy").(string),
				ArbitraryFlags:     is2ss(d.Get("swarm_opt").([]interface{})),
				ArbitraryJoinFlags: is2ss(d.Get("swarm_join_opt").([]interface{})),
				IsExperimental:     d.Get("swarm_experimental").(bool),
			},
		}

		exists, err := client.Exists(h.Name)
		if err != nil {
			return fmt.Errorf("Error checking if host exists: %s", err)
		}
		if exists {
			return mcnerror.ErrHostAlreadyExists{
				Name: h.Name,
			}
		}

		driverOpts := getDriverOpts(d, h.Driver.GetCreateFlags())

		if err := h.Driver.SetConfigFromFlags(driverOpts); err != nil {
			return fmt.Errorf("Error setting machine configuration from flags provided: %s", err)
		}

		if err := client.Create(h); err != nil {
			time.Sleep(2 * time.Second)

			vBoxLog := ""
			if h.DriverName == "virtualbox" {
				vBoxLog = filepath.Join(client.Path, "machines", h.Name, h.Name, "Logs", "VBox.log")
			}

			return crashreport.CrashError{
				Cause:       err,
				Command:     "Create",
				Context:     "client.performCreate",
				DriverName:  h.DriverName,
				LogFilePath: vBoxLog,
			}
		}

		if err := client.Save(h); err != nil {
			return fmt.Errorf("Error attempting to save store: %s", err)
		}

		d.Set("ssh_username", h.Driver.GetSSHUsername())
		d.Set("ssh_keypath", h.Driver.GetSSHKeyPath())
		machineState, err := h.Driver.GetState()
		if err != nil {
			return fmt.Errorf("Error attempting to retrieve state: %s", err)
		}
		switch machineState {
		case state.Timeout:
			return fmt.Errorf("Machine is in timeout state")
		case state.Error:
			return fmt.Errorf("Machine is in error state")
		}
		switch d.Get("state").(string) {
		case "running":
			switch machineState {
			case state.Paused, state.Saved, state.Stopped, state.Stopping:
				if err = h.Start(); err != nil {
					return fmt.Errorf("Error while attempting to start machine: %s", err)
				}
			}
		case "stopped":
			switch machineState {
			case state.Running, state.Starting:
				if err = h.Stop(); err != nil {
					return fmt.Errorf("Error while attempting to stop machine: %s", err)
				}
			}
		}
		machineState, err = h.Driver.GetState()
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
		}
		d.Set("state", strings.ToLower(machineState.String()))

		d.SetId(name)

		return nil
	}
}

func tlsPath(d *schema.ResourceData, option, directory, defaultValue string) string {
	ret := d.Get(option).(string)
	if len(ret) > 0 {
		return ret
	}
	return filepath.Join(directory, defaultValue)
}
