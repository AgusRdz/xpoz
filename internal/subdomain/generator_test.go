package subdomain

import (
	"context"
	"strings"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"a3f9c", true},
		{"bright-river", true},
		{"storm", true},
		{"ab", true},
		{"a-b", true},
		{"abc123", true},
		{"", false},
		{"a", false},                     // too short (< 2 chars)
		{"-foo", false},                  // leading hyphen
		{"foo-", false},                  // trailing hyphen
		{"FOO", false},                   // uppercase
		{"foo bar", false},               // space
		{"foo.bar", false},               // dot
		{"foo_bar", false},               // underscore
		{strings.Repeat("a", 64), false}, // too long (> 64 chars)
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			if got := IsValid(tc.input); got != tc.want {
				t.Errorf("IsValid(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestGenerateHex(t *testing.T) {
	lengths := []int{5, 6, 8, 0} // 0 should default to 5
	for _, l := range lengths {
		want := l
		if want <= 0 {
			want = 5
		}
		sub, err := Generate(StyleHex, l)
		if err != nil {
			t.Fatalf("Generate(hex, %d): %v", l, err)
		}
		if len(sub) != want {
			t.Errorf("Generate(hex, %d): got len %d, want %d", l, len(sub), want)
		}
		if !IsValid(sub) {
			t.Errorf("Generate(hex, %d): %q is not valid", l, sub)
		}
	}
}

func TestGenerateHex_unique(t *testing.T) {
	seen := make(map[string]bool, 100)
	for i := 0; i < 100; i++ {
		sub, err := generateHex(5)
		if err != nil {
			t.Fatalf("generateHex: %v", err)
		}
		seen[sub] = true
	}
	if len(seen) < 90 {
		t.Errorf("expected high uniqueness, got %d distinct values in 100 samples", len(seen))
	}
}

func TestGenerateSlug(t *testing.T) {
	sub, err := Generate(StyleSlug, 0)
	if err != nil {
		t.Fatalf("Generate(slug): %v", err)
	}
	if !strings.Contains(sub, "-") {
		t.Errorf("slug %q has no hyphen", sub)
	}
	if !IsValid(sub) {
		t.Errorf("slug %q is not valid", sub)
	}
	parts := strings.SplitN(sub, "-", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		t.Errorf("slug %q does not have two parts", sub)
	}
}

func TestGenerateWord(t *testing.T) {
	sub, err := Generate(StyleWord, 0)
	if err != nil {
		t.Fatalf("Generate(word): %v", err)
	}
	if !IsValid(sub) {
		t.Errorf("word %q is not valid", sub)
	}
	if strings.Contains(sub, "-") {
		t.Errorf("word %q contains hyphen", sub)
	}
}

func TestGenerate_unknownStyle(t *testing.T) {
	_, err := Generate("unknown", 5)
	if err == nil {
		t.Error("expected error for unknown style, got nil")
	}
}

func TestGenerateUnique_firstAttemptFree(t *testing.T) {
	ctx := context.Background()
	taken := func(_ context.Context, _ string) bool { return false }

	sub, err := GenerateUnique(ctx, StyleHex, 5, taken)
	if err != nil {
		t.Fatalf("GenerateUnique: %v", err)
	}
	if !IsValid(sub) {
		t.Errorf("GenerateUnique: %q is not valid", sub)
	}
}

func TestGenerateUnique_retriesOnCollision(t *testing.T) {
	ctx := context.Background()
	calls := 0
	// reject first 3 attempts
	taken := func(_ context.Context, _ string) bool {
		calls++
		return calls <= 3
	}

	sub, err := GenerateUnique(ctx, StyleHex, 5, taken)
	if err != nil {
		t.Fatalf("GenerateUnique: %v", err)
	}
	if calls < 4 {
		t.Errorf("expected at least 4 taken calls, got %d", calls)
	}
	if !IsValid(sub) {
		t.Errorf("GenerateUnique: %q is not valid", sub)
	}
}

func TestGenerateUnique_exhausted(t *testing.T) {
	ctx := context.Background()
	taken := func(_ context.Context, _ string) bool { return true }

	_, err := GenerateUnique(ctx, StyleHex, 5, taken)
	if err == nil {
		t.Error("expected error when all attempts taken, got nil")
	}
}

func TestWordLists_valid(t *testing.T) {
	for _, w := range adjectives {
		if w == "" {
			t.Errorf("empty adjective in list")
		}
	}
	for _, w := range nouns {
		if w == "" {
			t.Errorf("empty noun in list")
		}
	}
	for _, w := range words {
		if !IsValid(w) {
			t.Errorf("word %q fails IsValid", w)
		}
	}
}
