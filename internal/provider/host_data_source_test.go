package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestHostDataSource(t *testing.T) {
	t.Run("Schema", func(t *testing.T) {
		ctx := context.Background()
		schemaRequest := datasource.SchemaRequest{}
		schemaResponse := &datasource.SchemaResponse{}

		NewHostDataSource().Schema(ctx, schemaRequest, schemaResponse)

		if schemaResponse.Diagnostics.HasError() {
			t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
		}

		// Check that required attributes are present
		if _, exists := schemaResponse.Schema.Attributes["host"]; !exists {
			t.Fatal("Expected host attribute to be present in schema")
		}

		if _, exists := schemaResponse.Schema.Attributes["path"]; !exists {
			t.Fatal("Expected path attribute to be present in schema")
		}

		if _, exists := schemaResponse.Schema.Attributes["rendered"]; !exists {
			t.Fatal("Expected rendered attribute to be present in schema")
		}

		if _, exists := schemaResponse.Schema.Attributes["host_map"]; !exists {
			t.Fatal("Expected host_map attribute to be present in schema")
		}

		if _, exists := schemaResponse.Schema.Attributes["id"]; !exists {
			t.Fatal("Expected id attribute to be present in schema")
		}
	})

	t.Run("Metadata", func(t *testing.T) {
		ctx := context.Background()
		metadataRequest := datasource.MetadataRequest{
			ProviderTypeName: "sshconfig",
		}
		metadataResponse := &datasource.MetadataResponse{}

		NewHostDataSource().Metadata(ctx, metadataRequest, metadataResponse)

		if metadataResponse.TypeName != "sshconfig_host" {
			t.Fatalf("Expected type name to be 'sshconfig_host', got %s", metadataResponse.TypeName)
		}
	})
}

// TestAccHostDataSource tests the data source with actual Terraform configurations
func TestAccHostDataSource(t *testing.T) {
	// Create a temporary SSH config file for testing
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "ssh_config")

	configContent := `Host example.com
    HostName 192.168.1.100
    User testuser
    Port 2222
    IdentityFile ~/.ssh/id_rsa

Host *.example.org
    User admin
    ProxyJump bastion.example.org
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test SSH config file: %v", err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccHostDataSourceConfig(configFile, "example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.sshconfig_host.test", "host", "example.com"),
					resource.TestCheckResourceAttr("data.sshconfig_host.test", "path", configFile),
					resource.TestCheckResourceAttrSet("data.sshconfig_host.test", "rendered"),
					resource.TestCheckResourceAttrSet("data.sshconfig_host.test", "id"),
					resource.TestCheckResourceAttr("data.sshconfig_host.test", "host_map.HostName", "192.168.1.100"),
					resource.TestCheckResourceAttr("data.sshconfig_host.test", "host_map.User", "testuser"),
					resource.TestCheckResourceAttr("data.sshconfig_host.test", "host_map.Port", "2222"),
				),
			},
		},
	})
}

func testAccHostDataSourceConfig(path, host string) string {
	return fmt.Sprintf(`
data "sshconfig_host" "test" {
  path = %[1]q
  host = %[2]q
}
`, path, host)
}
