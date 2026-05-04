package controller

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
)

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name  string
		ports []container.Port
		want  string
	}{
		{"empty", nil, ""},
		{"single exposed port no binding", []container.Port{{PrivatePort: 80, Type: "tcp"}}, "80/tcp"},
		{"single port with public binding", []container.Port{{PrivatePort: 80, PublicPort: 8080, Type: "tcp"}}, "8080->80/tcp"},
		{"duplicate ports deduplicated", []container.Port{
			{PrivatePort: 80, Type: "tcp"},
			{PrivatePort: 80, Type: "tcp"},
		}, "80/tcp"},
		{"multiple distinct ports", []container.Port{
			{PrivatePort: 80, Type: "tcp"},
			{PrivatePort: 443, Type: "tcp"},
		}, "80/tcp, 443/tcp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPorts(tt.ports)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatCreated(t *testing.T) {
	t.Run("invalid string returned as-is", func(t *testing.T) {
		got := formatCreated("not-a-timestamp")
		assert.Equal(t, "not-a-timestamp", got)
	})

	t.Run("30 seconds ago", func(t *testing.T) {
		ts := time.Now().Add(-30 * time.Second).Format(time.RFC3339Nano)
		got := formatCreated(ts)
		assert.Equal(t, "30s ago", got)
	})

	t.Run("5 minutes ago", func(t *testing.T) {
		ts := time.Now().Add(-5 * time.Minute).Format(time.RFC3339Nano)
		got := formatCreated(ts)
		assert.Equal(t, "5m ago", got)
	})

	t.Run("3 hours ago", func(t *testing.T) {
		ts := time.Now().Add(-3 * time.Hour).Format(time.RFC3339Nano)
		got := formatCreated(ts)
		assert.Equal(t, "3h ago", got)
	})

	t.Run("2 days ago", func(t *testing.T) {
		ts := time.Now().Add(-48 * time.Hour).Format(time.RFC3339Nano)
		got := formatCreated(ts)
		assert.Equal(t, fmt.Sprintf("%dd ago", 48/24), got)
	})
}
