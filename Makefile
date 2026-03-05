# PorkLauncher Build System
# Two artifacts: bootstrap (small loader) + payload (full Wails app)

# ─── Configuration ───

VERSION        ?= dev
GOOS           ?= $(shell go env GOOS)
GOARCH         ?= $(shell go env GOARCH)

# Bootstrap build-time variables (override via env or make args)
UPDATE_SERVER_URL       ?=
ED25519_PUBLIC_KEY_HEX  ?=
FALLBACK_DOWNLOAD_URL   ?= https://github.com/Serezjjja/PorkLauncher/releases

# Payload build-time variables (injected via ldflags)
AZURIOM_BASE_URL    ?=
SESSION_SERVICE_URL ?=
PATCH_DOMAIN        ?=
PATCH_API_URL       ?=
GAME_PATCHES_URL    ?=
SERVER_LOGO_URL     ?=
SERVER_BANNER_URL   ?=
SERVER_IP           ?=
SERVER_NAME         ?=
JRE_MANIFEST_URL    ?=
BUTLER_BASE_URL     ?=
DISCORD_APP_ID      ?=

# ─── Derived ───

BS_EXT  := $(if $(filter windows,$(GOOS)),.exe,)
OUT_DIR := build/bin

BS_LDFLAGS := -s -w \
  -X 'HyLauncher/internal/bootstrap.UpdateServerURL=$(UPDATE_SERVER_URL)' \
  -X 'HyLauncher/internal/bootstrap.Ed25519PublicKeyHex=$(ED25519_PUBLIC_KEY_HEX)' \
  -X 'HyLauncher/internal/bootstrap.FallbackDownloadURL=$(FALLBACK_DOWNLOAD_URL)'

PAYLOAD_LDFLAGS := \
  -X 'HyLauncher/internal/config.AzuriomBaseURL=$(AZURIOM_BASE_URL)' \
  -X 'HyLauncher/internal/config.SessionServiceURL=$(SESSION_SERVICE_URL)' \
  -X 'HyLauncher/internal/config.PatchDomain=$(PATCH_DOMAIN)' \
  -X 'HyLauncher/internal/config.PatchAPIURL=$(PATCH_API_URL)' \
  -X 'HyLauncher/internal/config.GamePatchesURL=$(GAME_PATCHES_URL)' \
  -X 'HyLauncher/internal/config.ServerLogoURL=$(SERVER_LOGO_URL)' \
  -X 'HyLauncher/internal/config.ServerBannerURL=$(SERVER_BANNER_URL)' \
  -X 'HyLauncher/internal/config.ServerIP=$(SERVER_IP)' \
  -X 'HyLauncher/internal/config.ServerName=$(SERVER_NAME)' \
  -X 'HyLauncher/internal/config.JREManifestURL=$(JRE_MANIFEST_URL)' \
  -X 'HyLauncher/internal/config.ButlerBaseURL=$(BUTLER_BASE_URL)' \
  -X 'HyLauncher/internal/config.DiscordAppID=$(DISCORD_APP_ID)'

# ─── Targets ───

.PHONY: all bootstrap payload update-helper clean generate-keys sign-metadata help

all: bootstrap payload update-helper  ## Build everything

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Bootstrap (small, clean binary) ───

bootstrap:  ## Build bootstrap loader (~4-6 MB)
	@mkdir -p $(OUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="$(BS_LDFLAGS)" \
		-o $(OUT_DIR)/PorkLauncher$(BS_EXT) \
		./cmd/bootstrap
	@echo "Bootstrap built: $(OUT_DIR)/PorkLauncher$(BS_EXT)"
	@ls -lh $(OUT_DIR)/PorkLauncher$(BS_EXT)

# ─── Payload (Wails app — requires wails CLI) ───

payload:  ## Build payload Wails app (requires wails CLI)
	wails build -platform $(GOOS)/$(GOARCH) -o PorkLand -ldflags "$(PAYLOAD_LDFLAGS)"
	@echo "Payload built in $(OUT_DIR)/"

# ─── Update helper ───

update-helper:  ## Build update-helper binary
	@mkdir -p $(OUT_DIR)
	go build -ldflags="$(PAYLOAD_LDFLAGS)" \
		-o $(OUT_DIR)/update-helper$(BS_EXT) \
		./cmd/update-helper
	@echo "Update helper built: $(OUT_DIR)/update-helper$(BS_EXT)"

# ─── Crypto tools ───

generate-keys:  ## Generate Ed25519 keypair for payload signing
	go run ./tools/sign-payload -generate
	@echo ""
	@echo "Add the public key hex to GitHub Secrets as ED25519_PUBLIC_KEY_HEX"
	@echo "Add the private key hex to GitHub Secrets as ED25519_PRIVATE_KEY"

sign-metadata:  ## Sign metadata.json (requires private.key in cwd)
	@test -f metadata.json || (echo "ERROR: metadata.json not found" && exit 1)
	@test -f private.key || (echo "ERROR: private.key not found" && exit 1)
	go run ./tools/sign-payload -sign -key private.key -input metadata.json
	@echo "Signed: metadata.json.sig"

verify-metadata:  ## Verify metadata.json signature (requires public.key)
	@test -f metadata.json || (echo "ERROR: metadata.json not found" && exit 1)
	@test -f metadata.json.sig || (echo "ERROR: metadata.json.sig not found" && exit 1)
	@test -f public.key || (echo "ERROR: public.key not found" && exit 1)
	go run ./tools/sign-payload -verify -pubkey public.key -input metadata.json -sig metadata.json.sig

# ─── Cleanup ───

clean:  ## Remove build artifacts
	rm -rf $(OUT_DIR)
	rm -rf release/
	@echo "Cleaned build artifacts"
