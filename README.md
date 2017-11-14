# terraform-provider-dockermachine
Docker machine provider for Terraform

[![Go Report Card](https://goreportcard.com/badge/github.com/gstruct/terraform-provider-dockermachine)](https://goreportcard.com/report/github.com/gstruct/terraform-provider-dockermachine) [![Build Status](https://travis-ci.org/gstruct/terraform-provider-dockermachine.svg?branch=master)](https://travis-ci.org/gstruct/terraform-provider-dockermachine)

## Requisites

* [Terraform](https://www.terraform.io/)

**Note**: docker-machine is not required as its library is embedded into the provider.

## Install
```
$ go get github.com/gstruct/terraform-provider-dockermachine
```

## Usage

This provider makes available to Terraform all the docker-machine drivers as resources named "dockermachine\_\<drivername\>".  
All the creation flags of each driver (common or specific) are available as attributes of the resource, with dash characters ("-") replaced by underlines ("_").  
Furthermore, the following computed attributes are available:

* **address**: IP address of the docker machine
* **docker\_url**: URL of the docker daemon
* **docker\_version**: version of the docker daemon
* **ssh\_hostname**: SSH hostname
* **ssh\_keypath**: SSH private key path
* **ssh\_port**: SSH port
* **ssh\_username**: SSH username

Finally the state of the machine can be set using the attribute "state", either "running" or "stopped". Upon refresh, state will contain the actual state of the machine, lowercased.

Currently, any change to resource attributes, except for the "state" attribute, will trigger a destroy-create cycle.

The following parameters can be set at provider level:

* **debug**: boolean, enables docker-machine debug output in Terraform log
* **storage_path**: set default storage path for docker-machine
* **certs_directory**: set default path for docker-machine certs directory

### Example

```
resource "dockermachine_virtualbox" "node" {
    count = 2
    name = "${format("node-%02d", count.index+1)}"
    virtualbox_cpu_count = 2
    virtualbox_memory = 1024
    
    provisioner "remote-exec" {
        inline = [
            "touch /tmp/this_is_a_test",
        ]
        connection {
            type        = "ssh"
            host        = "${self.ssh_hostname}"
            port        = "${self.ssh_port}"
            user        = "${self.ssh_username}"
            private_key = "${file("${self.ssh_keypath}")}"
        }
    }
}
```
