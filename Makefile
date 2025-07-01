# Makefile for Terraform SSH Config Provider

.PHONY: build test docs clean install lint fmt

# Build the provider
build:
	go build -v .

# Run tests
test:
	go test -v ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v ./internal/provider/ -timeout 120m

# Generate documentation
docs:
	~/go/bin/tfplugindocs generate -provider-name sshconfig

# Install the provider locally for development
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/petems/sshconfig/0.1.0/linux_amd64
	cp terraform-provider-sshconfig ~/.terraform.d/plugins/registry.terraform.io/petems/sshconfig/0.1.0/linux_amd64/

# Clean build artifacts
clean:
	rm -f terraform-provider-sshconfig
	rm -rf docs/

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...
	@if command -v terraform >/dev/null 2>&1; then \
		terraform fmt -recursive ./examples/; \
	else \
		echo "terraform not found, skipping terraform fmt"; \
	fi

# Generate all (docs, format)
generate: fmt docs

# Development workflow
dev: fmt test build

# Full CI workflow
ci: fmt test testacc build docs

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the provider binary"
	@echo "  test     - Run unit tests"
	@echo "  testacc  - Run acceptance tests"
	@echo "  docs     - Generate documentation"
	@echo "  install  - Install provider locally"
	@echo "  clean    - Clean build artifacts"
	@echo "  lint     - Run linter"
	@echo "  fmt      - Format code"
	@echo "  generate - Generate docs and format"
	@echo "  dev      - Development workflow"
	@echo "  ci       - Full CI workflow"
	@echo "  help     - Show this help"
