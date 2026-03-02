package bootstrap

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// VerifyMetadataSignature verifies the Ed25519 signature of metadata.json.
// rawMetadata is the raw bytes of metadata.json as downloaded.
// signature is the raw Ed25519 signature (64 bytes).
func VerifyMetadataSignature(pubKey ed25519.PublicKey, rawMetadata, signature []byte) error {
	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("invalid signature size: got %d, want %d", len(signature), ed25519.SignatureSize)
	}

	if !ed25519.Verify(pubKey, rawMetadata, signature) {
		return fmt.Errorf("metadata signature verification failed: signature does not match")
	}

	return nil
}

// VerifyFileSHA256 checks that a file on disk matches the expected SHA256 hash.
// expectedHex should be a lowercase hex-encoded SHA256 hash (64 chars).
func VerifyFileSHA256(filePath string, expectedHex string) error {
	expectedHex = strings.ToLower(strings.TrimSpace(expectedHex))
	if len(expectedHex) != 64 {
		return fmt.Errorf("invalid SHA256 hash length: got %d chars, want 64", len(expectedHex))
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file for hash: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("hash file: %w", err)
	}

	actualHex := hex.EncodeToString(h.Sum(nil))
	if actualHex != expectedHex {
		return fmt.Errorf("SHA256 mismatch: expected %s, got %s", expectedHex, actualHex)
	}

	return nil
}
