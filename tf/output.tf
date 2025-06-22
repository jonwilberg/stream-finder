output "elasticsearch_node_ip_address" {
  description = "The public IPv4 address of the Droplet"
  value       = digitalocean_droplet.elasticsearch_node.ipv4_address
}