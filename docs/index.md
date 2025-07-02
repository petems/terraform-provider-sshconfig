---
page_title: "Provider: SSH Config"
description: |-
  The SSH Config provider allows you to read SSH configuration files and extract host configurations for use in Terraform configurations.
---

# SSH Config Provider

The SSH Config provider allows you to read SSH configuration files and extract host-specific configurations. This is useful for integrating existing SSH infrastructure configurations into your Terraform workflows.

## Example Usage

```terraform
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
  path = "/etc/ssh/ssh_config"  # Optional, defaults to /etc/ssh/ssh_config
}

# Access the rendered configuration
output "ssh_config" {
  value = data.sshconfig_host.example.rendered
}

# Access specific configuration values
output "hostname" {
  value = data.sshconfig_host.example.host_map["HostName"]
}

output "user" {
  value = data.sshconfig_host.example.host_map["User"]
}
```

## Schema

### Optional

- `path` (String) The path to the SSH config file to read. Defaults to `/etc/ssh/ssh_config`.

## Data Sources

- [sshconfig_host](./data-sources/host.md)

## SSH Config File Format

The provider reads standard SSH client configuration files. Here's an example:

```
Host example.com
    HostName 192.168.1.100
    User myuser
    Port 2222
    IdentityFile ~/.ssh/id_rsa

Host *.development
    User admin
    ProxyJump bastion.example.com
```

## Use Cases

- **Infrastructure Discovery**: Extract SSH connection details for dynamic infrastructure
- **Configuration Management**: Use existing SSH configurations in Terraform workflows  
- **Migration Planning**: Analyze existing SSH setups before infrastructure changes
- **Automation**: Integrate SSH configurations with other Terraform resources

## Limitations

- This provider is read-only and does not modify SSH configuration files
- The provider requires read access to the SSH configuration file
- Host pattern matching follows standard SSH configuration rules