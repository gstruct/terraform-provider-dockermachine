
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
