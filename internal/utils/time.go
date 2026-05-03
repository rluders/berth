package utils

import (
	"fmt"
	"time"
)

// FormatAge converts a Unix timestamp to a human-readable age string (e.g. "2h", "3d").
func FormatAge(unixTime int64) string {
	d := time.Since(time.Unix(unixTime, 0))
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// FormatBytes formats a byte count as a human-readable string (e.g. "512KB", "1.2GB").
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	val := float64(b) / float64(div)
	uc := "KMGTPE"[exp]
	if val == float64(uint64(val)) {
		return fmt.Sprintf("%d%cB", uint64(val), uc)
	}
	return fmt.Sprintf("%.1f%cB", val, uc)
}
