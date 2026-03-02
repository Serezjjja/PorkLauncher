// sign-payload is a CLI tool for Ed25519 key generation and metadata signing.
//
// Usage:
//
//	# Generate a new Ed25519 keypair
//	go run ./tools/sign-payload -generate
//
//	# Sign metadata.json with a private key
//	go run ./tools/sign-payload -sign -key private.key -input metadata.json
//
//	# Verify a signature
//	go run ./tools/sign-payload -verify -pubkey public.key -input metadata.json -sig metadata.json.sig
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	generateCmd := flag.Bool("generate", false, "Generate a new Ed25519 keypair")
	signCmd := flag.Bool("sign", false, "Sign a file")
	verifyCmd := flag.Bool("verify", false, "Verify a signature")

	keyPath := flag.String("key", "", "Path to private key file (hex-encoded)")
	pubkeyPath := flag.String("pubkey", "", "Path to public key file (hex-encoded)")
	inputPath := flag.String("input", "", "Path to input file to sign/verify")
	sigPath := flag.String("sig", "", "Path to signature file (for verify; for sign, auto-generated as input+'.sig')")

	flag.Parse()

	switch {
	case *generateCmd:
		doGenerate()
	case *signCmd:
		doSign(*keyPath, *inputPath)
	case *verifyCmd:
		doVerify(*pubkeyPath, *inputPath, *sigPath)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func doGenerate() {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fatalf("generate keypair: %v", err)
	}

	pubHex := hex.EncodeToString(pubKey)
	privHex := hex.EncodeToString(privKey)

	// Write files
	if err := os.WriteFile("private.key", []byte(privHex), 0600); err != nil {
		fatalf("write private key: %v", err)
	}
	if err := os.WriteFile("public.key", []byte(pubHex), 0644); err != nil {
		fatalf("write public key: %v", err)
	}

	fmt.Println("Ed25519 keypair generated successfully!")
	fmt.Println()
	fmt.Printf("Private key: private.key (KEEP SECRET!)\n")
	fmt.Printf("Public key:  public.key\n")
	fmt.Println()
	fmt.Printf("Public key hex (for -ldflags):\n")
	fmt.Printf("  -X 'HyLauncher/internal/bootstrap.Ed25519PublicKeyHex=%s'\n", pubHex)
	fmt.Println()
	fmt.Println("IMPORTANT: Store private.key securely. Add it to GitHub Secrets as ED25519_PRIVATE_KEY.")
	fmt.Println("           NEVER commit private.key to the repository!")
}

func doSign(keyPath, inputPath string) {
	if keyPath == "" || inputPath == "" {
		fatalf("usage: -sign -key <private.key> -input <metadata.json>")
	}

	privKeyHex, err := os.ReadFile(keyPath)
	if err != nil {
		fatalf("read private key: %v", err)
	}

	// Trim any whitespace/newlines that may have been introduced
	privKeyHexStr := strings.TrimSpace(string(privKeyHex))

	privKeyBytes, err := hex.DecodeString(privKeyHexStr)
	if err != nil {
		fatalf("decode private key hex: %v", err)
	}

	if len(privKeyBytes) != ed25519.PrivateKeySize {
		fatalf("invalid private key size: got %d bytes (%d hex chars), want %d bytes (%d hex chars)\n"+
			"Please regenerate keys with: go run ./tools/sign-payload -generate",
			len(privKeyBytes), len(privKeyHexStr), ed25519.PrivateKeySize, ed25519.PrivateKeySize*2)
	}

	privKey := ed25519.PrivateKey(privKeyBytes)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fatalf("read input file: %v", err)
	}

	signature := ed25519.Sign(privKey, data)

	sigOutputPath := inputPath + ".sig"
	if err := os.WriteFile(sigOutputPath, signature, 0644); err != nil {
		fatalf("write signature: %v", err)
	}

	fmt.Printf("Signed %s -> %s\n", inputPath, sigOutputPath)
	fmt.Printf("Signature (hex): %s\n", hex.EncodeToString(signature))
}

func doVerify(pubkeyPath, inputPath, sigPath string) {
	if pubkeyPath == "" || inputPath == "" || sigPath == "" {
		fatalf("usage: -verify -pubkey <public.key> -input <metadata.json> -sig <metadata.json.sig>")
	}

	pubKeyHex, err := os.ReadFile(pubkeyPath)
	if err != nil {
		fatalf("read public key: %v", err)
	}

	pubKeyHexStr := strings.TrimSpace(string(pubKeyHex))

	pubKeyBytes, err := hex.DecodeString(pubKeyHexStr)
	if err != nil {
		fatalf("decode public key hex: %v", err)
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		fatalf("invalid public key size: got %d bytes (%d hex chars), want %d bytes (%d hex chars)\n"+
			"Please regenerate keys with: go run ./tools/sign-payload -generate",
			len(pubKeyBytes), len(pubKeyHexStr), ed25519.PublicKeySize, ed25519.PublicKeySize*2)
	}

	pubKey := ed25519.PublicKey(pubKeyBytes)

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fatalf("read input file: %v", err)
	}

	sig, err := os.ReadFile(sigPath)
	if err != nil {
		fatalf("read signature file: %v", err)
	}

	if len(sig) != ed25519.SignatureSize {
		fatalf("invalid signature size: got %d bytes, want %d", len(sig), ed25519.SignatureSize)
	}

	if ed25519.Verify(pubKey, data, sig) {
		fmt.Println("Signature is VALID")
	} else {
		fmt.Println("Signature is INVALID")
		os.Exit(1)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
