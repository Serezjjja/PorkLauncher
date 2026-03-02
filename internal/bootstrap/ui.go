package bootstrap

import (
	"fmt"
	"os"
	"strings"
)

// UI handles minimal user-facing output for the bootstrap process.
// It uses stderr to avoid interfering with stdout if payload uses it.
type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) Status(msg string) {
	fmt.Fprintf(os.Stderr, "[PorkLauncher] %s\n", msg)
}

func (u *UI) Progress(downloaded, total int64) {
	if total > 0 {
		pct := float64(downloaded) / float64(total) * 100
		bar := progressBar(pct, 30)
		fmt.Fprintf(os.Stderr, "\r[PorkLauncher] Downloading: %s %.1f%% (%s / %s)",
			bar, pct, formatBytes(downloaded), formatBytes(total))
	} else {
		fmt.Fprintf(os.Stderr, "\r[PorkLauncher] Downloading: %s", formatBytes(downloaded))
	}
}

func (u *UI) ProgressDone() {
	fmt.Fprintln(os.Stderr)
}

func (u *UI) Error(msg string) {
	fmt.Fprintf(os.Stderr, "[PorkLauncher] ERROR: %s\n", msg)
}

func (u *UI) FallbackMessage() {
	fmt.Fprintf(os.Stderr, "\n"+
		"[PorkLauncher] Failed to load the launcher automatically.\n"+
		"[PorkLauncher] Please download the latest version manually:\n"+
		"[PorkLauncher]   %s\n\n",
		FallbackDownloadURL)
}

func progressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", width-filled) + "]"
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMG"[exp])
}
