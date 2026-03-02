# Bootstrap + Payload Architecture

## Overview

PorkLauncher uses a two-stage architecture:

1. **Bootstrap** (~6 MB) — A minimal loader with zero external dependencies.
   It contains only the update server URL and an Ed25519 public key.
   No game logic, no auth tokens, no frontend assets.

2. **Payload** (~15-30 MB) — The full Wails application (Go backend + React frontend).
   Downloaded, cryptographically verified, and launched by the bootstrap at runtime.

```
User runs PorkLauncher.exe (bootstrap)
  │
  ├─ Fetch metadata.json + metadata.json.sig
  ├─ Verify Ed25519 signature
  ├─ Check if cached payload matches latest version
  │   ├─ YES → launch cached payload
  │   └─ NO  → download payload ZIP
  │            ├─ Verify SHA256
  │            ├─ Extract to ~/.hylauncher/payload/<version>/
  │            └─ Launch HyLauncher.exe from extracted dir
  └─ On any error → show fallback URL for manual download
```

## Directory Structure

```
cmd/
├── bootstrap/main.go      ← Bootstrap entry point (go build)
├── payload/main.go         ← Payload entry point (copy of root main.go)
├── update-helper/          ← Update helper (unchanged)
└── verify/                 ← File verifier (unchanged)

internal/
├── bootstrap/              ← Bootstrap-only code
│   ├── config.go           ← UpdateServerURL, Ed25519PublicKeyHex (ldflags)
│   ├── fetch.go            ← HTTP: metadata + signature + payload download
│   ├── verify.go           ← Ed25519 signature + SHA256 verification
│   ├── extract.go          ← ZIP extraction with zip-slip protection
│   ├── launch.go           ← os/exec: start payload process
│   └── ui.go               ← Minimal stderr progress output
├── app/                    ← Payload business logic (unchanged)
├── config/                 ← Payload config (unchanged)
└── ...                     ← All other payload packages (unchanged)

tools/
└── sign-payload/main.go   ← CLI for key generation + signing

main.go                    ← Root entry point for Wails (payload, used by wails build)
```

## Setup Guide

### 1. Generate Ed25519 Keypair

```bash
make generate-keys
# or
go run ./tools/sign-payload -generate
```

This creates:
- `private.key` — **KEEP SECRET!** Add to GitHub Secrets as `ED25519_PRIVATE_KEY`
- `public.key` — Safe to share. The hex value goes into `ED25519_PUBLIC_KEY_HEX`

### 2. Configure GitHub Secrets

Add these new secrets to your repository:

| Secret | Value | Purpose |
|--------|-------|---------|
| `UPDATE_SERVER_URL` | URL where metadata.json is hosted | Bootstrap fetches from here |
| `ED25519_PUBLIC_KEY_HEX` | 64-char hex from `public.key` | Embedded in bootstrap binary |
| `ED25519_PRIVATE_KEY` | 128-char hex from `private.key` | Used in CI to sign metadata |
| `FALLBACK_DOWNLOAD_URL` | GitHub releases URL | Shown on error |

The existing secrets (AZURIOM_BASE_URL, etc.) remain unchanged — they're only used for the payload build.

### 3. Build Locally

```bash
# Build bootstrap only (no Wails needed)
make bootstrap \
  UPDATE_SERVER_URL=https://your-server.com/updates \
  ED25519_PUBLIC_KEY_HEX=$(cat public.key)

# Build payload (requires Wails CLI)
make payload

# Build everything
make all
```

### 4. CI/CD (GitHub Actions)

Both workflows have been updated:

- **build.yml** — Builds bootstrap + payload on every push/PR
- **manual-release.yml** — Creates a release with:
  - `PorkLauncher-bootstrap-*.exe` — Bootstrap binaries
  - `payload-*.zip` — Signed payload archives
  - `metadata.json` + `metadata.json.sig` — Signed metadata
  - `PorkLauncher-*-portable.exe` / `.dmg` — Direct download (legacy)
  - `version.json` — Legacy version manifest (backward compatibility)

## How Updates Work

### Publishing a New Version

1. **Trigger** the Manual Release workflow with a version number (e.g., `1.2.0`)
2. CI builds payload for all platforms → creates `payload-{os}-{arch}.zip`
3. CI generates `metadata.json` with SHA256 hashes and download URLs
4. CI signs `metadata.json` with the Ed25519 private key → `metadata.json.sig`
5. All artifacts are uploaded to the GitHub Release

### Update Server

The bootstrap fetches two files from `UPDATE_SERVER_URL`:

```
{UPDATE_SERVER_URL}/metadata.json
{UPDATE_SERVER_URL}/metadata.json.sig
```

**Option A: GitHub Releases (simplest)**
Set `UPDATE_SERVER_URL` to point to the latest release assets. You can use a redirect
service or GitHub Pages to serve the latest metadata.json.

**Option B: Static file hosting**
Upload `metadata.json` and `metadata.json.sig` to any HTTPS server (S3, Cloudflare R2,
your own server, etc.). Update them after each release.

**Option C: GitHub Pages**
Create a branch `gh-pages` with `metadata.json` + `metadata.json.sig`. Update after
each release via CI.

### metadata.json Format

```json
{
  "version": "1.2.0",
  "payload": {
    "windows/amd64": {
      "url": "https://github.com/.../payload-windows-amd64.zip",
      "sha256": "abc123...",
      "size": 45000000
    },
    "darwin/arm64": {
      "url": "https://github.com/.../payload-darwin-arm64.zip",
      "sha256": "def456...",
      "size": 42000000
    },
    "darwin/amd64": {
      "url": "https://github.com/.../payload-darwin-amd64.zip",
      "sha256": "789abc...",
      "size": 40000000
    }
  },
  "min_bootstrap_version": "1.0.0"
}
```

### Security Model

1. **Ed25519 public key** is compiled into the bootstrap binary at build time
2. `metadata.json.sig` is a 64-byte Ed25519 signature of the raw metadata bytes
3. Bootstrap verifies: `ed25519.Verify(pubKey, rawMetadataBytes, signature)`
4. After signature check, bootstrap verifies SHA256 of the downloaded payload ZIP
5. **Key rotation**: requires releasing a new bootstrap binary with the new public key

### Caching

Payloads are cached at:
- **Windows**: `%LOCALAPPDATA%\HyLauncher\payload\<version>\`
- **macOS**: `~/Library/Application Support/HyLauncher/payload/<version>/`
- **Linux**: `~/.hylauncher/payload/<version>/`

The current version is stored in `payload/current_version`. Old versions are
automatically cleaned up.

## Manual Signing

If you need to sign metadata outside of CI:

```bash
# Sign
go run ./tools/sign-payload -sign -key private.key -input metadata.json

# Verify
go run ./tools/sign-payload -verify -pubkey public.key -input metadata.json -sig metadata.json.sig
```

## Artifact Summary

| Artifact | Size | Contains |
|----------|------|----------|
| `PorkLauncher.exe` (bootstrap) | ~6 MB | Loader + pubkey + URL |
| `payload-windows-amd64.zip` | ~15-30 MB | HyLauncher.exe + update-helper.exe |
| `payload-darwin-arm64.zip` | ~15-30 MB | HyLauncher.app bundle |
| `metadata.json` | ~500 B | Version, URLs, SHA256 hashes |
| `metadata.json.sig` | 64 B | Ed25519 signature |
