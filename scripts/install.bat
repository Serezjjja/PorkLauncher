@echo off
:: PorkLauncher install wrapper for Windows
:: Double-click to install PorkLauncher

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1"
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [!] Installation failed. Press any key to exit.
    pause >nul
    exit /b %ERRORLEVEL%
)
