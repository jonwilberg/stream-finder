variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}
variable "ssh_key_name" {
  description = "SSH key name"
  type        = string
}
variable "ssh_key_path" {
  description = "Path to the SSH key"
  type        = string
}
variable "elasticsearch_password" {
  description = "Elasticsearch password"
  type        = string
  sensitive   = true
}