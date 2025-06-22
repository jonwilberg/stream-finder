terraform {
  required_version = ">= 1.0.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

provider "google" {
  project = "stream-finder"
  region  = "europe-north1"
}

data "digitalocean_ssh_key" "default" {
  name = var.ssh_key_name
}

data "digitalocean_project" "stream_finder" {
  name = "stream-finder"
}

resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "default" {
  name     = "${random_id.default.hex}-terraform-remote-backend"
  project  = var.gcp_project_id
  location = "europe-north1"

  force_destroy               = false
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }
}

resource "local_file" "default" {
  file_permission = "0644"
  filename        = "${path.module}/backend.tf"

  content = <<-EOT
terraform {
  backend "gcs" {
    bucket = "${google_storage_bucket.default.name}"
  }
}
EOT
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

resource "digitalocean_project_resources" "elasticsearch_node" {
  project = data.digitalocean_project.stream_finder.id
  resources = [
    digitalocean_droplet.elasticsearch_node.urn
  ]
}