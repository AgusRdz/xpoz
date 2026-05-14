// Package subdomain generates and validates tunnel subdomain names.
package subdomain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
)

// Style controls the format of generated subdomains.
type Style string

const (
	StyleHex  Style = "hex"  // "a3f9c" — 5 hex chars, 1M+ combinations
	StyleSlug Style = "slug" // "bright-river" — adjective-noun
	StyleWord Style = "word" // "storm" — single curated word
)

// Taken reports whether a subdomain is already in use.
// The caller supplies this function to avoid a package-level dependency on store.
type Taken func(ctx context.Context, sub string) bool

var validSubdomain = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}[a-z0-9]$`)

// IsValid reports whether s is a safe, DNS-compatible subdomain component.
func IsValid(s string) bool {
	return validSubdomain.MatchString(s)
}

// Generate creates a random subdomain in the given style.
// length is used only for StyleHex; other styles have fixed lengths.
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

// GenerateUnique generates a subdomain not already reported as taken.
// It retries up to 5 times before returning an error.
func GenerateUnique(ctx context.Context, style Style, length int, taken Taken) (string, error) {
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		sub, err := Generate(style, length)
		if err != nil {
			return "", err
		}
		if !taken(ctx, sub) {
			return sub, nil
		}
	}
	return "", fmt.Errorf("subdomain: could not generate unique value after %d attempts", maxAttempts)
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
	ai, err := randIndex(len(adjectives))
	if err != nil {
		return "", fmt.Errorf("subdomain: generating slug: %w", err)
	}
	ni, err := randIndex(len(nouns))
	if err != nil {
		return "", fmt.Errorf("subdomain: generating slug: %w", err)
	}
	return adjectives[ai] + "-" + nouns[ni], nil
}

func generateWord() (string, error) {
	i, err := randIndex(len(words))
	if err != nil {
		return "", fmt.Errorf("subdomain: generating word: %w", err)
	}
	return words[i], nil
}

func randIndex(n int) (int, error) {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}
