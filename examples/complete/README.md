# Complete SSH Config Provider Example

This example demonstrates all the features of the SSH Config provider, including:

- Reading SSH configurations for specific hosts
- Working with SSH host patterns and wildcards
- Extracting individual configuration values
- Using SSH config data with other Terraform resources
- Creating formatted SSH connection strings

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Apply the configuration:
   ```bash
   terraform apply
   ```

3. View the outputs:
   ```bash
   terraform output
   ```

## What This Example Demonstrates

### SSH Config File Creation
The example creates a test SSH config file with various host patterns:
- Specific hosts (`example.com`, `bastion.example.com`, `github.com`)
- Wildcard patterns (`*.development`, `production-*`)
- Different SSH configuration options (ports, users, identity files, proxy jumps, etc.)

### Data Source Usage
Shows how to read SSH configurations for different types of hosts:
- Direct host matches
- Wildcard pattern matches
- Hosts with various SSH options

### Output Examples
The example produces several outputs:
- **Full rendered configurations**: Complete SSH config blocks as strings
- **Connection details**: Structured data with individual SSH parameters
- **SSH connection strings**: Ready-to-use SSH command lines

### Integration with Other Resources
Demonstrates how to use SSH config data with:
- `null_resource` with SSH connections
- Dynamic SSH connection strings
- Conditional resource creation

## Expected Outputs

After running `terraform apply`, you'll see outputs like:

```
connection_details = {
  "example" = {
    "hostname" = "192.168.1.100"
    "key_file" = "~/.ssh/id_rsa"
    "port" = "2222"
    "user" = "myuser"
  }
  "development" = {
    "hostname" = "app1.development"
    "key_file" = "~/.ssh/dev_key"
    "proxy_jump" = "bastion.example.com"
    "user" = "admin"
  }
  # ... more hosts
}

ssh_connection_strings = {
  "example" = "ssh myuser@192.168.1.100 -p 2222 -i ~/.ssh/id_rsa"
  "development" = "ssh admin@app1.development -i ~/.ssh/dev_key"
  "bastion" = "ssh ubuntu@10.0.1.5 -i ~/.ssh/bastion_key"
  # ... more connection strings
}
```

## Notes

- The `null_resource` example is disabled by default (count = 0) to avoid requiring actual SSH connectivity
- To enable the SSH connection example, change `count = 0` to `count = 1` in the `null_resource`
- The example uses the `local` provider to create a test SSH config file
- All SSH configuration parsing follows standard OpenSSH client configuration rules