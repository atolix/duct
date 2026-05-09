package util

import (
	"fmt"
	"os"
	"strings"
)

func Humanize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func Shorten(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if rest, ok := strings.CutPrefix(path, home); ok {
		return "~" + rest
	}

	return path
}
