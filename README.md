# Terraform SSH Config Provider

A Terraform provider for reading SSH configuration files, built with the modern [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).

## Features

- **Modern Plugin Framework**: Built with Terraform Plugin Framework v1.15.0 for better performance and developer experience
- **Comprehensive Testing**: Includes both unit tests and acceptance tests
- **Rich Data Access**: Provides both rendered string output and structured map access to SSH configuration values
- **Flexible Configuration**: Supports custom SSH config file paths or uses system defaults
- **Auto-generated Documentation**: Complete documentation generated from schema definitions

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for development)

## Using the Provider

### Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    sshconfig = {
      source = "petems/sshconfig"
    }
  }
}
```

### Example Usage

```hcl
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

## Data Sources

### `sshconfig_host`

Reads SSH configuration for a specific host from an SSH config file.

#### Arguments

- `host` (Required) - The host you want to lookup in the SSH config file
- `path` (Optional) - The path to the SSH config file. Defaults to `/etc/ssh/ssh_config`

#### Attributes

- `id` - Unique identifier for the data source (format: `path:host`)
- `rendered` - The complete host configuration rendered as a multi-line string
- `host_map` - A map containing all configuration keys and values for the host

#### Example SSH Config File

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

## Documentation

Complete documentation is available in the `docs/` directory and is automatically generated from the provider schema:

- [Provider Documentation](docs/index.md) - Complete provider overview and usage
- [sshconfig_host Data Source](docs/data-sources/host.md) - Detailed data source documentation

### Generating Documentation

The documentation is automatically generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```bash
# Install the documentation tool
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# Generate documentation
make docs
```

The documentation generation:
- Extracts schema information directly from the provider code
- Generates markdown files with proper formatting
- Includes examples and usage patterns
- Validates that documentation stays in sync with code changes

## Development

### Building the Provider

```bash
go build -v .
```

### Running Tests

#### Unit Tests
```bash
go test -v ./internal/provider/
```

#### Acceptance Tests
```bash
TF_ACC=1 go test -v ./internal/provider/
```

### Makefile Targets

A comprehensive Makefile is provided for development workflows:

```bash
make help          # Show all available targets
make build         # Build the provider
make test          # Run unit tests
make testacc       # Run acceptance tests
make docs          # Generate documentation
make fmt           # Format code and examples
make lint          # Run linter
make dev           # Development workflow (fmt, lint, test, build)
make ci            # Full CI workflow (includes docs generation)
```

### Local Development

To use a locally built provider:

1. Build the provider:
   ```bash
   make build
   ```

2. Install locally:
   ```bash
   make install
   ```

3. Create a `.terraformrc` file in your home directory:
   ```hcl
   provider_installation {
     dev_overrides {
       "petems/sshconfig" = "/path/to/your/provider/binary"
     }
     direct {}
   }
   ```

## Migration from SDK v2

This provider has been upgraded from Terraform Plugin SDK v2 to the modern Terraform Plugin Framework. Key improvements include:

### What's New

- **Better Error Handling**: More descriptive error messages with proper diagnostic support
- **Improved Type Safety**: Stronger typing with the Plugin Framework's type system
- **Enhanced Testing**: Comprehensive unit and acceptance tests using the latest testing patterns
- **Modern Go Patterns**: Uses idiomatic Go patterns and interfaces instead of declarative structs
- **Future-Proof**: Built on the recommended framework for new Terraform providers
- **Auto-generated Documentation**: Schema-driven documentation generation

### Breaking Changes

- **Provider Address**: The provider source has been updated to use the modern registry format
- **Go Version**: Now requires Go 1.23+ for development
- **Terraform Version**: Requires Terraform 1.0+ (maintains compatibility with 0.12+ at runtime)

### Migration Guide

If you're upgrading from an older version:

1. Update your provider source in `terraform` blocks:
   ```hcl
   # Old
   terraform {
     required_providers {
       sshconfig = {
         source = "github.com/petems/terraform-provider-sshconfig"
       }
     }
   }

   # New
   terraform {
     required_providers {
       sshconfig = {
         source = "petems/sshconfig"
       }
     }
   }
   ```

2. The data source interface remains the same, so existing configurations should work without changes.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `go test -v ./...`
6. Generate documentation: `make docs`
7. Submit a pull request

### Code Quality

This project maintains high code quality standards:

- All code must pass `go vet` and `golangci-lint`
- Test coverage should be maintained or improved
- Documentation must be generated and up-to-date
- Follow the [Terraform Plugin Framework best practices](https://developer.hashicorp.com/terraform/plugin/framework)

### Documentation Requirements

- All provider schema changes must include updated documentation
- Examples should be provided for new features
- Documentation is automatically generated and validated in CI/CD
- Templates in `templates/` directory control documentation structure

## Examples

Comprehensive examples are available in the `examples/` directory:

- [Basic Usage](examples/data-sources/sshconfig_host/) - Simple data source usage
- [Complete Example](examples/complete/) - Comprehensive feature demonstration

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- Uses [go-sshconfig](https://github.com/petems/go-sshconfig) for SSH config parsing
- Documentation generated with [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)
- Follows [HashiCorp's provider development guidelines](https://developer.hashicorp.com/terraform/plugin)