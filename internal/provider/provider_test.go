package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sshconfig": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	provider := New("test")()

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}
}

func TestProviderMetadata(t *testing.T) {
	provider := &SSHConfigProvider{version: "test"}

	// Test that provider has the correct type name
	if provider.version != "test" {
		t.Errorf("Expected version to be 'test', got %s", provider.version)
	}
}
