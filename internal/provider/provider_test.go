package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
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
	prov := New("test")()

	if prov == nil {
		t.Fatal("Expected provider to be created")
	}
}

func TestProviderMetadata(t *testing.T) {
	prov := &SSHConfigProvider{version: "test"}

	// Test that provider has the correct type name
	if prov.version != "test" {
		t.Errorf("Expected version to be 'test', got %s", prov.version)
	}
}

func TestProviderMetadataRequest(t *testing.T) {
	ctx := context.Background()
	prov := &SSHConfigProvider{version: "1.0.0"}

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	prov.Metadata(ctx, req, resp)

	if resp.TypeName != "sshconfig" {
		t.Errorf("Expected TypeName to be 'sshconfig', got %s", resp.TypeName)
	}

	if resp.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got %s", resp.Version)
	}
}
