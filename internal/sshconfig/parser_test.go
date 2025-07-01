package sshconfig

import (
	"strings"
	"testing"
)

func TestTokenizeLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasError bool
	}{
		{
			name:     "simple tokens",
			input:    "Host example.com",
			expected: []string{"Host", "example.com"},
		},
		{
			name:     "multiple tokens",
			input:    "HostName 192.168.1.100",
			expected: []string{"HostName", "192.168.1.100"},
		},
		{
			name:     "quoted string",
			input:    `ProxyCommand "ssh -W %h:%p bastion"`,
			expected: []string{"ProxyCommand", "ssh -W %h:%p bastion"},
		},
		{
			name:     "single quoted string",
			input:    `User 'test user'`,
			expected: []string{"User", "test user"},
		},
		{
			name:     "multiple spaces",
			input:    "Host     example.com    test.com",
			expected: []string{"Host", "example.com", "test.com"},
		},
		{
			name:     "tabs",
			input:    "Host\texample.com\ttest.com",
			expected: []string{"Host", "example.com", "test.com"},
		},
		{
			name:     "unterminated quote",
			input:    `Host "unterminated`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tokenizeLine(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(result))
				return
			}

			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("token %d: expected %q, got %q", i, tt.expected[i], token)
				}
			}
		})
	}
}

func TestProperCaseKeyword(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hostname", "HostName"},
		{"port", "Port"},
		{"user", "User"},
		{"identityfile", "IdentityFile"},
		{"stricthostkeychecking", "StrictHostKeyChecking"},
		{"proxyjump", "ProxyJump"},
		{"forwardagent", "ForwardAgent"},
		{"unknownkeyword", "Unknownkeyword"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := properCaseKeyword(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		pattern  string
		hostname string
		expected bool
	}{
		{"example.com", "example.com", true},
		{"example.com", "test.com", false},
		{"*.example.com", "test.example.com", true},
		{"*.example.com", "example.com", false},
		{"*", "anything", true},
		{"test*", "test123", true},
		{"test*", "123test", false},
		{"!example.com", "example.com", false},
		{"!example.com", "test.com", false}, // negation doesn't match, so returns false
		{"*.dev", "app.dev", true},
		{"*.dev", "app.prod", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.hostname, func(t *testing.T) {
			result := matchPattern(tt.pattern, tt.hostname)
			if result != tt.expected {
				t.Errorf("pattern %q with hostname %q: expected %v, got %v",
					tt.pattern, tt.hostname, tt.expected, result)
			}
		})
	}
}

func TestHostMatches(t *testing.T) {
	host := &Host{
		Patterns: []string{"*.example.com", "test.com"},
		Options:  make(map[string]string),
	}

	tests := []struct {
		hostname string
		expected bool
	}{
		{"app.example.com", true},
		{"test.com", true},
		{"example.com", false},
		{"other.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			result := host.Matches(tt.hostname)
			if result != tt.expected {
				t.Errorf("hostname %q: expected %v, got %v", tt.hostname, tt.expected, result)
			}
		})
	}
}

func TestParseSimpleConfig(t *testing.T) {
	configText := `
# This is a comment
Host example.com
    HostName 192.168.1.100
    User myuser
    Port 2222

Host *.dev
    User admin
    ProxyJump bastion.example.com
`

	config, err := Parse(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if len(config.Hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(config.Hosts))
	}

	// Test first host
	host1 := config.Hosts[0]
	if len(host1.Patterns) != 1 || host1.Patterns[0] != "example.com" {
		t.Errorf("first host patterns: expected [example.com], got %v", host1.Patterns)
	}

	expectedOptions1 := map[string]string{
		"HostName": "192.168.1.100",
		"User":     "myuser",
		"Port":     "2222",
	}

	for key, expected := range expectedOptions1 {
		if actual, exists := host1.Options[key]; !exists {
			t.Errorf("first host missing option %q", key)
		} else if actual != expected {
			t.Errorf("first host option %q: expected %q, got %q", key, expected, actual)
		}
	}

	// Test second host
	host2 := config.Hosts[1]
	if len(host2.Patterns) != 1 || host2.Patterns[0] != "*.dev" {
		t.Errorf("second host patterns: expected [*.dev], got %v", host2.Patterns)
	}

	expectedOptions2 := map[string]string{
		"User":      "admin",
		"ProxyJump": "bastion.example.com",
	}

	for key, expected := range expectedOptions2 {
		if actual, exists := host2.Options[key]; !exists {
			t.Errorf("second host missing option %q", key)
		} else if actual != expected {
			t.Errorf("second host option %q: expected %q, got %q", key, expected, actual)
		}
	}
}

func TestFindHost(t *testing.T) {
	configText := `
Host example.com
    User myuser
    Port 2222

Host *.dev
    User admin

Host *
    User default
`

	config, err := Parse(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	tests := []struct {
		hostname     string
		expectedUser string
		shouldFind   bool
	}{
		{"example.com", "myuser", true},
		{"app.dev", "admin", true},
		{"unknown.com", "default", true},
		{"test.example.com", "default", true},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			host := config.FindHost(tt.hostname)

			if !tt.shouldFind {
				if host != nil {
					t.Errorf("expected not to find host for %q, but found one", tt.hostname)
				}
				return
			}

			if host == nil {
				t.Errorf("expected to find host for %q, but didn't", tt.hostname)
				return
			}

			if user, exists := host.Options["User"]; !exists {
				t.Errorf("host for %q missing User option", tt.hostname)
			} else if user != tt.expectedUser {
				t.Errorf("host for %q: expected user %q, got %q", tt.hostname, tt.expectedUser, user)
			}
		})
	}
}

func TestGetMergedOptions(t *testing.T) {
	configText := `
Host example.com
    User myuser
    Port 2222

Host *.com
    User defaultuser
    IdentityFile ~/.ssh/id_rsa

Host *
    Port 22
    IdentityFile ~/.ssh/default_key
`

	config, err := Parse(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	// Test that first match wins for each option
	merged := config.GetMergedOptions("example.com")

	expected := map[string]string{
		"User":         "myuser",        // from first host
		"Port":         "2222",          // from first host
		"IdentityFile": "~/.ssh/id_rsa", // from second host (*.com)
	}

	for key, expectedValue := range expected {
		if actualValue, exists := merged[key]; !exists {
			t.Errorf("missing option %q", key)
		} else if actualValue != expectedValue {
			t.Errorf("option %q: expected %q, got %q", key, expectedValue, actualValue)
		}
	}

	// Test with a hostname that matches *.com but not the specific example.com
	merged2 := config.GetMergedOptions("test.com")

	if user, exists := merged2["User"]; !exists {
		t.Errorf("missing User option for test.com")
	} else if user != "defaultuser" {
		t.Errorf("User for test.com: expected defaultuser, got %q", user)
	}

	if port, exists := merged2["Port"]; !exists {
		t.Errorf("missing Port option for test.com")
	} else if port != "22" {
		t.Errorf("Port for test.com: expected 22, got %q", port)
	}

	// Test with a hostname that only matches the wildcard *
	merged3 := config.GetMergedOptions("test.org")

	if port, exists := merged3["Port"]; !exists {
		t.Errorf("missing Port option for test.org")
	} else if port != "22" {
		t.Errorf("Port for test.org: expected 22, got %q", port)
	}

	if identityFile, exists := merged3["IdentityFile"]; !exists {
		t.Errorf("missing IdentityFile option for test.org")
	} else if identityFile != "~/.ssh/default_key" {
		t.Errorf("IdentityFile for test.org: expected ~/.ssh/default_key, got %q", identityFile)
	}

	// test.org should not have a User option since it doesn't match *.com or example.com
	if _, exists := merged3["User"]; exists {
		t.Errorf("test.org should not have User option, but it does")
	}
}

func TestHostString(t *testing.T) {
	host := &Host{
		Patterns: []string{"example.com", "*.dev"},
		Options: map[string]string{
			"HostName": "192.168.1.100",
			"User":     "myuser",
			"Port":     "2222",
		},
	}

	result := host.String()

	// Check that it starts with Host line
	if !strings.HasPrefix(result, "Host example.com *.dev\n") {
		t.Errorf("unexpected Host line in output: %s", result)
	}

	// Check that all options are present
	expectedOptions := []string{
		"HostName 192.168.1.100",
		"User myuser",
		"Port 2222",
	}

	for _, expected := range expectedOptions {
		if !strings.Contains(result, expected) {
			t.Errorf("missing expected option %q in output: %s", expected, result)
		}
	}
}

func TestParseMultipleHostPatterns(t *testing.T) {
	configText := `
Host example.com test.com *.dev
    User multihost
    Port 2222
`

	config, err := Parse(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if len(config.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(config.Hosts))
	}

	host := config.Hosts[0]
	expectedPatterns := []string{"example.com", "test.com", "*.dev"}

	if len(host.Patterns) != len(expectedPatterns) {
		t.Fatalf("expected %d patterns, got %d", len(expectedPatterns), len(host.Patterns))
	}

	for i, expected := range expectedPatterns {
		if host.Patterns[i] != expected {
			t.Errorf("pattern %d: expected %q, got %q", i, expected, host.Patterns[i])
		}
	}

	// Test that all patterns match correctly
	testHosts := []string{"example.com", "test.com", "app.dev"}
	for _, testHost := range testHosts {
		if !host.Matches(testHost) {
			t.Errorf("host %q should match but doesn't", testHost)
		}
	}
}

func TestParseEmptyAndCommentLines(t *testing.T) {
	configText := `
# This is a comment

Host example.com
    # Another comment
    User myuser
    
    Port 2222

# Final comment
`

	config, err := Parse(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if len(config.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(config.Hosts))
	}

	host := config.Hosts[0]
	expectedOptions := map[string]string{
		"User": "myuser",
		"Port": "2222",
	}

	for key, expected := range expectedOptions {
		if actual, exists := host.Options[key]; !exists {
			t.Errorf("missing option %q", key)
		} else if actual != expected {
			t.Errorf("option %q: expected %q, got %q", key, expected, actual)
		}
	}
}
