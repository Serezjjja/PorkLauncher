#!/usr/bin/env bash
# PorkLauncher install script — Linux & macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/Serezjjja/PorkLauncher/main/scripts/install.sh | bash
set -euo pipefail

REPO="Serezjjja/PorkLauncher"
APP_NAME="HyLauncher"
DISPLAY_NAME="PorkLauncher"

# ── Colors ───────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
info()  { echo -e "${GREEN}[+]${NC} $*"; }
warn()  { echo -e "${YELLOW}[!]${NC} $*"; }
error() { echo -e "${RED}[✗]${NC} $*"; exit 1; }

# ── Detect OS & arch ─────────────────────────────────────────────────────────
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Linux)  PLATFORM="linux" ;;
  Darwin) PLATFORM="macos" ;;
  *)      error "Unsupported OS: $OS" ;;
esac

case "$ARCH" in
  x86_64)         ARCH_TAG="x64"   ;;
  aarch64|arm64)  ARCH_TAG="arm64" ;;
  *)              error "Unsupported architecture: $ARCH" ;;
esac

# ── Dependencies ─────────────────────────────────────────────────────────────
command -v curl &>/dev/null || error "'curl' is required but not installed."

# ── Fetch latest release tag ─────────────────────────────────────────────────
info "Fetching latest $DISPLAY_NAME release..."
LATEST_TAG=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
[ -z "$LATEST_TAG" ] && error "Could not determine latest release. Check your internet connection."
info "Version: $LATEST_TAG"

BASE_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}"

# ── Linux ─────────────────────────────────────────────────────────────────────
if [ "$PLATFORM" = "linux" ]; then
  INSTALL_DIR="${XDG_BIN_HOME:-$HOME/.local/bin}"
  BINARY_PATH="$INSTALL_DIR/$APP_NAME"

  mkdir -p "$INSTALL_DIR"

  info "Downloading $DISPLAY_NAME for Linux ($ARCH_TAG)..."
  curl -fsSL --progress-bar \
    -o "$BINARY_PATH" \
    "${BASE_URL}/PorkLauncher-linux-${ARCH_TAG}"
  chmod +x "$BINARY_PATH"

  # ── .desktop entry ──────────────────────────────────────────────────────
  DESKTOP_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/applications"
  mkdir -p "$DESKTOP_DIR"
  cat > "$DESKTOP_DIR/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=$DISPLAY_NAME
Comment=Unofficial Hytale Launcher for free to play gamers
Exec=$BINARY_PATH
Icon=$APP_NAME
Categories=Game;Utility;
Keywords=hytale;launcher;gaming;mods;
Terminal=false
StartupNotify=true
EOF
  chmod +x "$DESKTOP_DIR/$APP_NAME.desktop"
  command -v update-desktop-database &>/dev/null \
    && update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true

  info "Installed    → $BINARY_PATH"
  info "Desktop entry → $DESKTOP_DIR/$APP_NAME.desktop"

  info "Launching $DISPLAY_NAME..."
  nohup "$BINARY_PATH" </dev/null >/dev/null 2>&1 &

# ── macOS ─────────────────────────────────────────────────────────────────────
elif [ "$PLATFORM" = "macos" ]; then
  TMP_DIR=$(mktemp -d)
  trap 'rm -rf "$TMP_DIR"' EXIT

  DMG_NAME="PorkLauncher-macos-${ARCH_TAG}.dmg"
  ZIP_NAME="PorkLauncher-macos-${ARCH_TAG}.zip"

  info "Downloading $DMG_NAME..."
  if curl -fsSL --progress-bar \
      -o "$TMP_DIR/$DMG_NAME" "${BASE_URL}/${DMG_NAME}" 2>/dev/null; then

    info "Mounting disk image..."
    MOUNT_POINT=$(hdiutil attach "$TMP_DIR/$DMG_NAME" -nobrowse -quiet \
      | tail -1 | awk '{print $NF}')

    info "Copying $APP_NAME.app to /Applications..."
    cp -R "$MOUNT_POINT/$APP_NAME.app" /Applications/
    hdiutil detach "$MOUNT_POINT" -quiet
  else
    warn "DMG not available — trying ZIP fallback..."
    curl -fsSL --progress-bar \
      -o "$TMP_DIR/$ZIP_NAME" "${BASE_URL}/${ZIP_NAME}" \
      || error "Failed to download $DISPLAY_NAME for macOS $ARCH_TAG."
    unzip -q "$TMP_DIR/$ZIP_NAME" -d "$TMP_DIR"
    cp -R "$TMP_DIR/$APP_NAME.app" /Applications/
  fi

  info "Installed → /Applications/$APP_NAME.app"

  info "Launching $DISPLAY_NAME..."
  open -a "$APP_NAME"
fi

echo ""
info "$DISPLAY_NAME installed successfully!"
