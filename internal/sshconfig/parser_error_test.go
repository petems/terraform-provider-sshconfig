package sshconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFileErrorHandling(t *testing.T) {
	t.Run("FileNotFound", func(t *testing.T) {
		_, err := ParseFile("/non/existent/file.conf")
		if err == nil {
			t.Error("Expected error when file does not exist")
		}

		if !strings.Contains(err.Error(), "failed to open file") {
			t.Errorf("Expected 'failed to open file' error, got: %v", err)
		}
	})

	t.Run("ValidFile", func(t *testing.T) {
		// Create a temporary file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "ssh_config")

		configContent := `Host example.com
    User testuser
    Port 2222
`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		config, err := ParseFile(configFile)
		if err != nil {
			t.Fatalf("ParseFile failed: %v", err)
		}

		if len(config.Hosts) != 1 {
			t.Errorf("Expected 1 host, got %d", len(config.Hosts))
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "empty_config")

		err := os.WriteFile(configFile, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		config, err := ParseFile(configFile)
		if err != nil {
			t.Fatalf("ParseFile failed: %v", err)
		}

		if len(config.Hosts) != 0 {
			t.Errorf("Expected 0 hosts for empty file, got %d", len(config.Hosts))
		}
	})
}

func TestParseErrorHandling(t *testing.T) {
	t.Run("InvalidHostDirective", func(t *testing.T) {
		configText := `Host
    User testuser
`

		_, err := Parse(strings.NewReader(configText))
		if err == nil {
			t.Error("Expected error for Host directive without patterns")
		}

		if !strings.Contains(err.Error(), "Host directive requires at least one pattern") {
			t.Errorf("Expected 'Host directive requires at least one pattern' error, got: %v", err)
		}
	})

	t.Run("UnterminatedQuote", func(t *testing.T) {
		configText := `Host "unterminated
    User testuser
`

		_, err := Parse(strings.NewReader(configText))
		if err == nil {
			t.Error("Expected error for unterminated quote")
		}

		if !strings.Contains(err.Error(), "unterminated quoted string") {
			t.Errorf("Expected 'unterminated quoted string' error, got: %v", err)
		}
	})

	t.Run("OnlyComments", func(t *testing.T) {
		configText := `# This is a comment
# Another comment
# Final comment
`

		config, err := Parse(strings.NewReader(configText))
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(config.Hosts) != 0 {
			t.Errorf("Expected 0 hosts for comment-only file, got %d", len(config.Hosts))
		}
	})

	t.Run("GlobalOptions", func(t *testing.T) {
		configText := `User globaluser
Port 2222

Host example.com
    User specificuser
`

		config, err := Parse(strings.NewReader(configText))
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if len(config.Hosts) != 2 {
			t.Errorf("Expected 2 hosts (global + specific), got %d", len(config.Hosts))
		}

		// Global options should be in the first host with * pattern
		globalHost := config.Hosts[0]
		if len(globalHost.Patterns) != 1 || globalHost.Patterns[0] != "*" {
			t.Errorf("Expected global host with '*' pattern, got %v", globalHost.Patterns)
		}

		if globalHost.Options["User"] != "globaluser" {
			t.Errorf("Expected global user 'globaluser', got %s", globalHost.Options["User"])
		}

		if globalHost.Options["Port"] != "2222" {
			t.Errorf("Expected global port '2222', got %s", globalHost.Options["Port"])
		}
	})
}

func TestMatchPatternEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		hostname string
		expected bool
	}{
		{
			name:     "empty pattern",
			pattern:  "",
			hostname: "example.com",
			expected: false,
		},
		{
			name:     "empty hostname",
			pattern:  "example.com",
			hostname: "",
			expected: false,
		},
		{
			name:     "both empty",
			pattern:  "",
			hostname: "",
			expected: true,
		},
		{
			name:     "complex wildcard",
			pattern:  "*-prod-*",
			hostname: "web-prod-01",
			expected: true,
		},
		{
			name:     "complex wildcard no match",
			pattern:  "*-prod-*",
			hostname: "web-dev-01",
			expected: false,
		},
		{
			name:     "question mark wildcard",
			pattern:  "web-??.example.com",
			hostname: "web-01.example.com",
			expected: true,
		},
		{
			name:     "question mark no match",
			pattern:  "web-??.example.com",
			hostname: "web-001.example.com",
			expected: false,
		},
		{
			name:     "negation exact match",
			pattern:  "!example.com",
			hostname: "example.com",
			expected: false,
		},
		{
			name:     "negation no match",
			pattern:  "!example.com",
			hostname: "test.com",
			expected: false,
		},
		{
			name:     "negation wildcard",
			pattern:  "!*.example.com",
			hostname: "test.example.com",
			expected: false,
		},
		{
			name:     "negation wildcard no match",
			pattern:  "!*.example.com",
			hostname: "test.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPattern(tt.pattern, tt.hostname)
			if result != tt.expected {
				t.Errorf("matchPattern(%q, %q): expected %v, got %v",
					tt.pattern, tt.hostname, tt.expected, result)
			}
		})
	}
}

func TestParserDepthLimit(t *testing.T) {
	parser := NewParser()

	// Test that parser has a reasonable depth limit
	if parser.maxDepth != 5 {
		t.Errorf("Expected maxDepth to be 5, got %d", parser.maxDepth)
	}

	// Test that includeDepth starts at 0
	if parser.includeDepth != 0 {
		t.Errorf("Expected includeDepth to be 0, got %d", parser.includeDepth)
	}
}

func TestGetMergedOptionsEdgeCases(t *testing.T) {
	t.Run("NoMatchingHosts", func(t *testing.T) {
		configText := `Host example.com
    User testuser
`

		config, err := Parse(strings.NewReader(configText))
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		merged := config.GetMergedOptions("nomatch.com")

		if len(merged) != 0 {
			t.Errorf("Expected empty options for non-matching host, got %d options", len(merged))
		}
	})

	t.Run("OverlappingPatterns", func(t *testing.T) {
		configText := `Host *.example.com
    User broaduser
    Port 2222

Host web.example.com
    User specificuser

Host *
    Port 22
    IdentityFile ~/.ssh/id_rsa
`

		config, err := Parse(strings.NewReader(configText))
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		merged := config.GetMergedOptions("web.example.com")

		// Should get User from the most specific match (web.example.com)
		if merged["User"] != "specificuser" {
			t.Errorf("Expected 'specificuser', got %s", merged["User"])
		}

		// Should get Port from the first match (*.example.com)
		if merged["Port"] != "2222" {
			t.Errorf("Expected '2222', got %s", merged["Port"])
		}

		// Should get IdentityFile from the wildcard match (*)
		if merged["IdentityFile"] != "~/.ssh/id_rsa" {
			t.Errorf("Expected '~/.ssh/id_rsa', got %s", merged["IdentityFile"])
		}
	})
}

func TestHostStringFormatting(t *testing.T) {
	t.Run("EmptyHost", func(t *testing.T) {
		host := &Host{
			Patterns: []string{},
			Options:  map[string]string{},
		}

		result := host.String()

		if !strings.HasPrefix(result, "Host \n") {
			t.Errorf("Expected 'Host \\n', got %q", result)
		}
	})

	t.Run("HostWithManyOptions", func(t *testing.T) {
		host := &Host{
			Patterns: []string{"example.com"},
			Options: map[string]string{
				"User":                  "testuser",
				"Port":                  "2222",
				"HostName":              "192.168.1.100",
				"IdentityFile":          "~/.ssh/id_rsa",
				"StrictHostKeyChecking": "no",
				"UserKnownHostsFile":    "/dev/null",
				"ForwardAgent":          "yes",
				"ProxyJump":             "bastion.example.com",
			},
		}

		result := host.String()

		// Should start with Host line
		if !strings.HasPrefix(result, "Host example.com\n") {
			t.Errorf("Expected to start with 'Host example.com\\n', got %q", result)
		}

		// All options should be present
		for key, value := range host.Options {
			expectedLine := fmt.Sprintf("    %s %s", key, value)
			if !strings.Contains(result, expectedLine) {
				t.Errorf("Missing option line %q in output", expectedLine)
			}
		}
	})
}
