package updater

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// verifySignature verifies the Ed25519 signature over data.
// If publicKey is empty, verification is skipped.
func verifySignature(data []byte, sigHex string) error {
	if publicKey == "" {
		return nil
	}
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return fmt.Errorf("decoding public key: %w", err)
	}
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return fmt.Errorf("decoding signature: %w", err)
	}
	if !ed25519.Verify(ed25519.PublicKey(pubBytes), data, sig) {
		return fmt.Errorf("signature verification failed — release may be tampered")
	}
	return nil
}

// verifyChecksum checks SHA256(data) == expected (hex string).
func verifyChecksum(data []byte, expected string) error {
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	if got != strings.ToLower(expected) {
		return fmt.Errorf("checksum mismatch: got %s, want %s", got[:16]+"...", expected[:16]+"...")
	}
	return nil
}

// findChecksum parses a checksums.txt file and returns the hash for filename.
func findChecksum(checksums []byte, filename string) (string, error) {
	for _, line := range strings.Split(string(checksums), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == filename {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("no checksum entry for %q in checksums.txt", filename)
}
