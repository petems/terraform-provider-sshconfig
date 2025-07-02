package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/petems/terraform-provider-sshconfig/internal/sshconfig"
)

func TestHostDataSourceErrorHandling(t *testing.T) {
	t.Run("ParserErrorHandling", func(t *testing.T) {
		// Test that the parser correctly handles file not found
		_, err := sshconfig.ParseFile("/non/existent/path")
		if err == nil {
			t.Error("Expected error when file does not exist")
		}
	})

	t.Run("HostNotFoundInConfig", func(t *testing.T) {
		// Create a temporary SSH config file without the host we're looking for
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "ssh_config")

		configContent := `Host other.com
    User otheruser
`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test SSH config file: %v", err)
		}

		// Test using our parser directly
		config, err := sshconfig.ParseFile(configFile)
		if err != nil {
			t.Fatalf("Parse should succeed: %v", err)
		}

		// Test that host is not found
		host := config.FindHost("nonexistent.com")
		if host != nil {
			t.Error("Expected nil host for non-existent host")
		}
	})

	t.Run("InvalidConfigSyntax", func(t *testing.T) {
		// Create a temporary SSH config file with invalid syntax
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "ssh_config")

		configContent := `Host "unterminated quote
    User test
`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test SSH config file: %v", err)
		}

		// Test that parser handles invalid config
		_, err = sshconfig.ParseFile(configFile)
		if err == nil {
			t.Error("Expected error when config has invalid syntax")
		}
	})
}

func TestAccHostDataSourceErrorHandling(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccHostDataSourceConfigNonExistentFile(),
				ExpectError: regexp.MustCompile("Unable to Parse SSH Config File"),
			},
			{
				Config:      testAccHostDataSourceConfigNonExistentHost(),
				ExpectError: regexp.MustCompile("Host Not Found"),
			},
		},
	})
}

func testAccHostDataSourceConfigNonExistentFile() string {
	return `
data "sshconfig_host" "test" {
  path = "/non/existent/file"
  host = "example.com"
}
`
}

func testAccHostDataSourceConfigNonExistentHost() string {
	// Create a temporary file for this test
	tempDir := os.TempDir()
	configFile := filepath.Join(tempDir, "test_ssh_config")

	configContent := `Host other.com
    User otheruser
`

	os.WriteFile(configFile, []byte(configContent), 0644)

	return fmt.Sprintf(`
data "sshconfig_host" "test" {
  path = %q
  host = "nonexistent.com"
}
`, configFile)
}
