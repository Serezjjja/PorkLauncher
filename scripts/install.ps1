#Requires -Version 5.1
# PorkLauncher install script — Windows
# Usage (PowerShell): iex (irm 'https://raw.githubusercontent.com/Serezjjja/PorkLauncher/main/scripts/install.ps1')
[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$ProgressPreference    = "SilentlyContinue"

$REPO        = "Serezjjja/PorkLauncher"
$APP_NAME    = "PorkLauncher"
$EXE_NAME    = "HyLauncher.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\$APP_NAME"

function Write-Step([string]$msg) { Write-Host "[+] $msg" -ForegroundColor Green  }
function Write-Warn([string]$msg) { Write-Host "[!] $msg" -ForegroundColor Yellow }
function Write-Err ([string]$msg) { Write-Host "[x] $msg" -ForegroundColor Red; exit 1 }

# ── Fetch latest release ──────────────────────────────────────────────────────
Write-Step "Fetching latest $APP_NAME release..."
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
} catch {
    Write-Err "Failed to fetch release info: $_"
}

$tag = $release.tag_name
if (-not $tag) { Write-Err "Could not determine latest release tag." }
Write-Step "Version: $tag"

$baseUrl      = "https://github.com/$REPO/releases/download/$tag"
$portableName = "PorkLauncher-windows-x64-portable.exe"
$downloadUrl  = "$baseUrl/$portableName"
$exePath      = "$INSTALL_DIR\$EXE_NAME"

# ── Install directory ─────────────────────────────────────────────────────────
New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null

# ── Download portable binary ──────────────────────────────────────────────────
Write-Step "Downloading $portableName..."
try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $exePath
} catch {
    Write-Err "Download failed: $_"
}

# ── Desktop shortcut ──────────────────────────────────────────────────────────
Write-Step "Creating desktop shortcut..."
$desktop      = [Environment]::GetFolderPath("Desktop")
$shortcutPath = "$desktop\$APP_NAME.lnk"

$shell              = New-Object -ComObject WScript.Shell
$shortcut           = $shell.CreateShortcut($shortcutPath)
$shortcut.TargetPath       = $exePath
$shortcut.WorkingDirectory = $INSTALL_DIR
$shortcut.Description      = "Unofficial Hytale Launcher for free to play gamers"
$shortcut.Save()

Write-Step "Installed    -> $exePath"
Write-Step "Shortcut     -> $shortcutPath"

# ── Launch ────────────────────────────────────────────────────────────────────
Write-Step "Launching $APP_NAME..."
Start-Process $exePath

Write-Host ""
Write-Step "$APP_NAME installed successfully!"
