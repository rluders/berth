package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatAge(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name     string
		unixTime int64
		want     string
	}{
		{"zero seconds", now, "0s"},
		{"30 seconds", now - 30, "30s"},
		{"59 seconds", now - 59, "59s"},
		{"exactly 60 seconds", now - 60, "1m"},
		{"5 minutes", now - 300, "5m"},
		{"59 minutes", now - 59*60, "59m"},
		{"exactly 1 hour", now - 3600, "1h"},
		{"3 hours", now - 3*3600, "3h"},
		{"23 hours", now - 23*3600, "23h"},
		{"exactly 1 day", now - 24*3600, "1d"},
		{"7 days", now - 7*24*3600, "7d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatAge(tt.unixTime)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		b    uint64
		want string
	}{
		{"zero", 0, "0B"},
		{"1 byte", 1, "1B"},
		{"1023 bytes", 1023, "1023B"},
		{"exactly 1 KB", 1024, "1KB"},
		{"1.5 KB", 1536, "1.5KB"},
		{"exactly 1 MB", 1024 * 1024, "1MB"},
		{"1.5 MB", 1024*1024 + 512*1024, "1.5MB"},
		{"exactly 1 GB", 1024 * 1024 * 1024, "1GB"},
		{"2 GB", 2 * 1024 * 1024 * 1024, "2GB"},
		{"1.2 GB", 1288490188, "1.2GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatBytes(tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
