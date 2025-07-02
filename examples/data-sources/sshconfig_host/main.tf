terraform {
  required_providers {
    sshconfig = {
      source = "petems/sshconfig"
    }
  }
}

# Read SSH configuration for a specific host
data "sshconfig_host" "example" {
  host = "example.com"
  path = "/etc/ssh/ssh_config"
}

# Output the rendered SSH configuration
output "ssh_config_rendered" {
  description = "The SSH configuration rendered as a string"
  value       = data.sshconfig_host.example.rendered
}

# Output specific SSH configuration values
output "ssh_hostname" {
  description = "The hostname for the SSH connection"
  value       = lookup(data.sshconfig_host.example.host_map, "HostName", "")
}

output "ssh_user" {
  description = "The user for the SSH connection"
  value       = lookup(data.sshconfig_host.example.host_map, "User", "")
}

output "ssh_port" {
  description = "The port for the SSH connection"
  value       = lookup(data.sshconfig_host.example.host_map, "Port", "22")
}

# Example using the SSH config data in other resources
locals {
  ssh_connection_string = "${lookup(data.sshconfig_host.example.host_map, "User", "root")}@${lookup(data.sshconfig_host.example.host_map, "HostName", data.sshconfig_host.example.host)}"
}

output "ssh_connection_string" {
  description = "A formatted SSH connection string"
  value       = local.ssh_connection_string
}