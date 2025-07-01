package sshconfig

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Config represents a parsed SSH configuration
type Config struct {
	Hosts []*Host
}

// Host represents a host block in SSH config
type Host struct {
	Patterns []string          // Host patterns (e.g., "*.example.com", "host1")
	Options  map[string]string // SSH options for this host
}

// Parser handles SSH config parsing
type Parser struct {
	includeDepth int
	maxDepth     int
}

// NewParser creates a new SSH config parser
func NewParser() *Parser {
	return &Parser{
		maxDepth: 5, // Prevent infinite recursion in includes
	}
}

// Parse parses an SSH config from a reader
func (p *Parser) Parse(r io.Reader) (*Config, error) {
	config := &Config{}
	scanner := bufio.NewScanner(r)

	var currentHost *Host
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line
		if err := p.parseLine(line, &currentHost, config, lineNum); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
	}

	// Add the last host if it exists
	if currentHost != nil {
		config.Hosts = append(config.Hosts, currentHost)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	return config, nil
}

// parseLine parses a single line of SSH config
func (p *Parser) parseLine(line string, currentHost **Host, config *Config, lineNum int) error {
	// Split the line into tokens, respecting quotes
	tokens, err := tokenizeLine(line)
	if err != nil {
		return fmt.Errorf("failed to tokenize line: %w", err)
	}

	if len(tokens) == 0 {
		return nil
	}

	keyword := strings.ToLower(tokens[0])

	switch keyword {
	case "host":
		// Save previous host if it exists
		if *currentHost != nil {
			config.Hosts = append(config.Hosts, *currentHost)
		}

		// Start new host block
		if len(tokens) < 2 {
			return fmt.Errorf("Host directive requires at least one pattern")
		}

		*currentHost = &Host{
			Patterns: tokens[1:],
			Options:  make(map[string]string),
		}

	case "match":
		// For simplicity, we'll treat Match blocks similar to Host blocks
		// In a full implementation, you'd parse the match criteria
		if *currentHost != nil {
			config.Hosts = append(config.Hosts, *currentHost)
		}

		*currentHost = &Host{
			Patterns: []string{"*"}, // Match blocks apply to all hosts by default
			Options:  make(map[string]string),
		}

	default:
		// This is a configuration option
		if *currentHost == nil {
			// Global options - create a default host block
			*currentHost = &Host{
				Patterns: []string{"*"},
				Options:  make(map[string]string),
			}
		}

		// Store the option
		if len(tokens) >= 2 {
			// Convert keyword back to proper case
			properKeyword := properCaseKeyword(keyword)
			value := strings.Join(tokens[1:], " ")
			(*currentHost).Options[properKeyword] = value
		}
	}

	return nil
}

// tokenizeLine splits a line into tokens, respecting quotes
func tokenizeLine(line string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		char := line[i]

		switch {
		case !inQuotes && (char == '"' || char == '\''):
			// Start of quoted string
			inQuotes = true
			quoteChar = char

		case inQuotes && char == quoteChar:
			// End of quoted string
			inQuotes = false
			quoteChar = 0

		case !inQuotes && (char == ' ' || char == '\t'):
			// Whitespace outside quotes - end current token
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		default:
			// Regular character
			current.WriteByte(char)
		}
	}

	// Add the last token
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unterminated quoted string")
	}

	return tokens, nil
}

// properCaseKeyword converts a lowercase keyword to proper case
func properCaseKeyword(keyword string) string {
	// Common SSH config keywords with proper casing
	keywordMap := map[string]string{
		"hostname":                         "HostName",
		"port":                             "Port",
		"user":                             "User",
		"identityfile":                     "IdentityFile",
		"identitiesonly":                   "IdentitiesOnly",
		"stricthostkeychecking":            "StrictHostKeyChecking",
		"userknownhostsfile":               "UserKnownHostsFile",
		"globalknownhostsfile":             "GlobalKnownHostsFile",
		"compression":                      "Compression",
		"compressionlevel":                 "CompressionLevel",
		"connectionattempts":               "ConnectionAttempts",
		"connecttimeout":                   "ConnectTimeout",
		"forwardagent":                     "ForwardAgent",
		"forwardx11":                       "ForwardX11",
		"forwardx11trusted":                "ForwardX11Trusted",
		"gatewayports":                     "GatewayPorts",
		"passwordauthentication":           "PasswordAuthentication",
		"pubkeyauthentication":             "PubkeyAuthentication",
		"rsaauthentication":                "RSAAuthentication",
		"batchmode":                        "BatchMode",
		"checkhostip":                      "CheckHostIP",
		"ciphers":                          "Ciphers",
		"clearallforwardings":              "ClearAllForwardings",
		"enablesshkeysign":                 "EnableSSHKeysign",
		"escapecchar":                      "EscapeChar",
		"exitonforwardfailure":             "ExitOnForwardFailure",
		"fallbacktoremotehost":             "FallBackToRsh",
		"hashknownhosts":                   "HashKnownHosts",
		"hostbasedauthentication":          "HostbasedAuthentication",
		"hostkeyalgorithms":                "HostKeyAlgorithms",
		"hostkeyalias":                     "HostKeyAlias",
		"kbdinteractiveauthentication":     "KbdInteractiveAuthentication",
		"localcommand":                     "LocalCommand",
		"localforward":                     "LocalForward",
		"loglevel":                         "LogLevel",
		"macs":                             "MACs",
		"nohostauthenticationforlocalhost": "NoHostAuthenticationForLocalhost",
		"numberofpasswordprompts":          "NumberOfPasswordPrompts",
		"permitlocalcommand":               "PermitLocalCommand",
		"preferredauthentications":         "PreferredAuthentications",
		"protocol":                         "Protocol",
		"proxycommand":                     "ProxyCommand",
		"proxyjump":                        "ProxyJump",
		"remoteforward":                    "RemoteForward",
		"rhostsrsaauthentication":          "RhostsRSAAuthentication",
		"sendenv":                          "SendEnv",
		"serveraliveinterval":              "ServerAliveInterval",
		"serveralivecountmax":              "ServerAliveCountMax",
		"tcpkeepalive":                     "TCPKeepAlive",
		"tunnel":                           "Tunnel",
		"tunneldevice":                     "TunnelDevice",
		"useprivilegedport":                "UsePrivilegedPort",
		"verifyhostkeydns":                 "VerifyHostKeyDNS",
		"visualhostkey":                    "VisualHostKey",
		"xauthlocation":                    "XAuthLocation",
		"addkeystoagent":                   "AddKeysToAgent",
		"include":                          "Include",
	}

	if proper, exists := keywordMap[keyword]; exists {
		return proper
	}

	// If not in our map, capitalize first letter of each word
	words := strings.Split(keyword, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

// FindHost finds the first host configuration that matches the given hostname
func (c *Config) FindHost(hostname string) *Host {
	for _, host := range c.Hosts {
		if host.Matches(hostname) {
			return host
		}
	}
	return nil
}

// Matches checks if this host configuration matches the given hostname
func (h *Host) Matches(hostname string) bool {
	for _, pattern := range h.Patterns {
		if matchPattern(pattern, hostname) {
			return true
		}
	}
	return false
}

// matchPattern implements SSH config pattern matching
func matchPattern(pattern, hostname string) bool {
	// Handle negation patterns (starting with !)
	if strings.HasPrefix(pattern, "!") {
		negPattern := pattern[1:]
		if negPattern == hostname {
			return false
		}
		if strings.Contains(negPattern, "*") || strings.Contains(negPattern, "?") {
			matched, _ := filepath.Match(negPattern, hostname)
			return !matched
		}
		// For negation patterns, if they don't match exactly, they don't apply (return false)
		return false
	}

	// Handle exact matches
	if pattern == hostname {
		return true
	}

	// Handle wildcard patterns
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		matched, _ := filepath.Match(pattern, hostname)
		return matched
	}

	return false
}

// String returns a string representation of the host configuration
func (h *Host) String() string {
	var result strings.Builder

	// Write host patterns
	result.WriteString("Host ")
	result.WriteString(strings.Join(h.Patterns, " "))
	result.WriteString("\n")

	// Write options
	for key, value := range h.Options {
		result.WriteString(fmt.Sprintf("    %s %s\n", key, value))
	}

	return result.String()
}

// GetMergedOptions returns the merged options for a specific hostname
// This applies SSH's precedence rules with specificity-based priority
func (c *Config) GetMergedOptions(hostname string) map[string]string {
	merged := make(map[string]string)
	optionSources := make(map[string]*Host) // Track which host provided each option

	// Collect all matching hosts with their specificity scores
	type hostWithScore struct {
		host  *Host
		score int
	}

	var matchingHosts []hostWithScore

	for _, host := range c.Hosts {
		if host.Matches(hostname) {
			score := calculateSpecificity(host, hostname)
			matchingHosts = append(matchingHosts, hostWithScore{host: host, score: score})
		}
	}

	// Process hosts in file order, but allow more specific hosts to override less specific ones
	for _, hostWithScore := range matchingHosts {
		host := hostWithScore.host
		currentScore := hostWithScore.score

		for key, value := range host.Options {
			existingSource, exists := optionSources[key]
			if !exists {
				// First time seeing this option
				merged[key] = value
				optionSources[key] = host
			} else {
				// Option already exists, check if current host is more specific
				existingScore := calculateSpecificity(existingSource, hostname)
				if currentScore > existingScore {
					merged[key] = value
					optionSources[key] = host
				}
			}
		}
	}

	return merged
}

// calculateSpecificity returns a specificity score for a host pattern matching a hostname
// Higher scores indicate more specific matches
func calculateSpecificity(host *Host, hostname string) int {
	maxScore := 0

	for _, pattern := range host.Patterns {
		score := 0

		if pattern == hostname {
			// Exact match is most specific
			score = 1000
		} else if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
			// Wildcard patterns get points based on literal character count
			literalChars := 0
			for _, char := range pattern {
				if char != '*' && char != '?' {
					literalChars++
				}
			}
			score = literalChars
		}

		if score > maxScore {
			maxScore = score
		}
	}

	return maxScore
}

// ParseFile parses an SSH config file from a file path
func ParseFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	parser := NewParser()
	return parser.Parse(file)
}

// Parse is a convenience function that creates a parser and parses the config
func Parse(r io.Reader) (*Config, error) {
	parser := NewParser()
	return parser.Parse(r)
}
