// Package subdomain generates and validates tunnel subdomain names.
package subdomain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
)

// Style controls the format of generated subdomains.
type Style string

const (
	StyleHex  Style = "hex"  // "a3f9c" — 5 hex chars, 1M+ combinations
	StyleSlug Style = "slug" // "blue-river" — adjective-noun
	StyleWord Style = "word" // "storm" — single word
)

// validSubdomain matches safe, DNS-compatible subdomain components.
var validSubdomain = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}[a-z0-9]$`)

// IsValid reports whether s is a valid subdomain component.
func IsValid(s string) bool {
	return validSubdomain.MatchString(s)
}

// Generate creates a random subdomain in the given style.
// length is only used for StyleHex; other styles have fixed lengths.
func Generate(style Style, length int) (string, error) {
	switch style {
	case StyleHex, "":
		return generateHex(length)
	case StyleSlug:
		return generateSlug()
	case StyleWord:
		return generateWord()
	default:
		return "", fmt.Errorf("subdomain: unknown style %q", style)
	}
}

func generateHex(length int) (string, error) {
	if length <= 0 {
		length = 5
	}
	b := make([]byte, (length+1)/2)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("subdomain: generating hex: %w", err)
	}
	return hex.EncodeToString(b)[:length], nil
}

func generateSlug() (string, error) {
	// TODO: implement adjective-noun generator in Phase 1
	return "", fmt.Errorf("subdomain: slug style not yet implemented")
}

func generateWord() (string, error) {
	// TODO: implement word generator in Phase 1
	return "", fmt.Errorf("subdomain: word style not yet implemented")
}
