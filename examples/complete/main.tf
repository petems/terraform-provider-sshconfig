terraform {
  required_providers {
    sshconfig = {
      source = "petems/sshconfig"
    }
  }
}

# Example SSH config file for demonstration
resource "local_file" "ssh_config" {
  filename = "${path.module}/test_ssh_config"
  content = <<-EOF
    Host example.com
        HostName 192.168.1.100
        User myuser
        Port 2222
        IdentityFile ~/.ssh/id_rsa
        ForwardAgent yes
        
    Host *.development
        User admin
        ProxyJump bastion.example.com
        IdentityFile ~/.ssh/dev_key
        
    Host bastion.example.com
        HostName 10.0.1.5
        User ubuntu
        IdentityFile ~/.ssh/bastion_key
        Port 22
        
    Host production-*
        User deploy
        IdentityFile ~/.ssh/prod_key
        StrictHostKeyChecking yes
        
    Host github.com
        HostName github.com
        User git
        IdentityFile ~/.ssh/github_key
        AddKeysToAgent yes
  EOF
}

# Read SSH configuration for a specific host
data "sshconfig_host" "example" {
  host = "example.com"
  path = local_file.ssh_config.filename
}

# Read configuration for a wildcard pattern host
data "sshconfig_host" "dev_server" {
  host = "app1.development"
  path = local_file.ssh_config.filename
}

# Read configuration for bastion host
data "sshconfig_host" "bastion" {
  host = "bastion.example.com"
  path = local_file.ssh_config.filename
}

# Read configuration for production host pattern
data "sshconfig_host" "production" {
  host = "production-web01"
  path = local_file.ssh_config.filename
}

# Read configuration for GitHub
data "sshconfig_host" "github" {
  host = "github.com"
  path = local_file.ssh_config.filename
}

# Output the rendered configurations
output "example_config" {
  description = "Full SSH configuration for example.com"
  value       = data.sshconfig_host.example.rendered
}

output "dev_config" {
  description = "Full SSH configuration for development servers"
  value       = data.sshconfig_host.dev_server.rendered
}

# Output specific configuration values
output "connection_details" {
  description = "SSH connection details for various hosts"
  value = {
    example = {
      hostname = data.sshconfig_host.example.host_map["HostName"]
      user     = data.sshconfig_host.example.host_map["User"]
      port     = lookup(data.sshconfig_host.example.host_map, "Port", "22")
      key_file = data.sshconfig_host.example.host_map["IdentityFile"]
    }
    
    development = {
      hostname    = lookup(data.sshconfig_host.dev_server.host_map, "HostName", "app1.development")
      user        = data.sshconfig_host.dev_server.host_map["User"]
      proxy_jump  = data.sshconfig_host.dev_server.host_map["ProxyJump"]
      key_file    = data.sshconfig_host.dev_server.host_map["IdentityFile"]
    }
    
    bastion = {
      hostname = data.sshconfig_host.bastion.host_map["HostName"]
      user     = data.sshconfig_host.bastion.host_map["User"]
      port     = data.sshconfig_host.bastion.host_map["Port"]
      key_file = data.sshconfig_host.bastion.host_map["IdentityFile"]
    }
    
    production = {
      user                    = data.sshconfig_host.production.host_map["User"]
      key_file               = data.sshconfig_host.production.host_map["IdentityFile"]
      strict_host_key_checking = data.sshconfig_host.production.host_map["StrictHostKeyChecking"]
    }
    
    github = {
      hostname        = data.sshconfig_host.github.host_map["HostName"]
      user           = data.sshconfig_host.github.host_map["User"]
      key_file       = data.sshconfig_host.github.host_map["IdentityFile"]
      add_keys_to_agent = data.sshconfig_host.github.host_map["AddKeysToAgent"]
    }
  }
}

# Demonstrate using SSH config data with other resources
resource "null_resource" "example_connection" {
  count = 0 # Set to 1 to enable this example
  
  connection {
    type        = "ssh"
    host        = data.sshconfig_host.example.host_map["HostName"]
    user        = data.sshconfig_host.example.host_map["User"]
    port        = lookup(data.sshconfig_host.example.host_map, "Port", 22)
    private_key = file(data.sshconfig_host.example.host_map["IdentityFile"])
  }

  provisioner "remote-exec" {
    inline = [
      "echo 'Connected to ${data.sshconfig_host.example.host_map["HostName"]} as ${data.sshconfig_host.example.host_map["User"]}'",
      "hostname",
      "whoami"
    ]
  }
}

# Create a formatted SSH connection string
locals {
  ssh_connections = {
    for host_key, host_data in {
      example     = data.sshconfig_host.example
      development = data.sshconfig_host.dev_server
      bastion     = data.sshconfig_host.bastion
      production  = data.sshconfig_host.production
      github      = data.sshconfig_host.github
    } : host_key => format(
      "ssh %s@%s%s%s",
      host_data.host_map["User"],
      lookup(host_data.host_map, "HostName", host_data.host),
      lookup(host_data.host_map, "Port", "22") != "22" ? " -p ${host_data.host_map["Port"]}" : "",
      can(host_data.host_map["IdentityFile"]) ? " -i ${host_data.host_map["IdentityFile"]}" : ""
    )
  }
}

output "ssh_connection_strings" {
  description = "Ready-to-use SSH connection strings"
  value       = local.ssh_connections
}

# Demonstrate error handling (commented out to avoid errors)
# data "sshconfig_host" "nonexistent" {
#   host = "nonexistent.example.com"
#   path = local_file.ssh_config.filename
# }