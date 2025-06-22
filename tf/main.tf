terraform {
  required_version = ">= 1.0.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "default" {
  name = var.ssh_key_name
}

resource "digitalocean_droplet" "elasticsearch_node" {
  name       = "elasticsearch-node"
  region     = "ams3"
  size       = "s-1vcpu-2gb"
  image      = "ubuntu-22-04-x64"
  ssh_keys   = [data.digitalocean_ssh_key.default.id]
  backups    = false
  monitoring = true

  connection {
    host        = self.ipv4_address
    type        = "ssh"
    user        = "root"
    private_key = file(var.ssh_key_path)
  }

  provisioner "file" {
    source      = "${path.module}/scripts/install_elasticsearch.sh"
    destination = "/tmp/install_elasticsearch.sh"
  }

  provisioner "file" {
    source      = "${path.module}/config/elasticsearch.yml"
    destination = "/etc/elasticsearch/elasticsearch.yml"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/install_elasticsearch.sh",
      "ELASTICSEARCH_PASSWORD='${var.elasticsearch_password}' /tmp/install_elasticsearch.sh",
    ]
  }
}