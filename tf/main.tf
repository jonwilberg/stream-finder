terraform {
  required_providers {
    google = { source = "hashicorp/google", version = ">= 4.56.0" }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_storage_bucket" "terraform_state" {
  name     = "tf-remote-backend-463006"
  location = var.region

  force_destroy               = false
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }
}

resource "local_file" "backend_config" {
  file_permission = "0644"
  filename        = "${path.module}/backend.tf"

  content = <<-EOT
  terraform {
    backend "gcs" {
      bucket = "${google_storage_bucket.terraform_state.name}"
    }
  }
  EOT
}

resource "google_project_service" "firestore_api" {
  project                    = var.project_id
  service                    = "firestore.googleapis.com"
  disable_dependent_services = true
}

resource "google_firestore_database" "firestore_database" {
  project     = var.project_id
  name        = "primary-prod"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"

  depends_on = [ google_project_service.firestore_api ]
}

