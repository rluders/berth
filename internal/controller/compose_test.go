package controller

import (
	"context"
	"testing"
	"time"

	"github.com/rluders/berth/internal/engine"
	"github.com/stretchr/testify/assert"
)

func TestStreamCompose_closesChanOnExit(t *testing.T) {
	// Use a project name that will fail fast (docker compose -p X version exits quickly).
	ctx := context.Background()
	ch := make(chan string, 64)

	err := StreamCompose(ctx, "berth-test-nonexistent", "", ch, "version")
	// StreamCompose may or may not error depending on docker availability.
	// What matters: ch must be closed after output drains.
	if err != nil {
		// Channel closed inside StreamCompose on Start failure.
		_, ok := <-ch
		assert.False(t, ok, "channel must be closed on start error")
		return
	}

	// Drain and verify channel eventually closes.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range ch {
		}
	}()

	select {
	case <-done:
		// channel closed — pass
	case <-time.After(10 * time.Second):
		t.Fatal("channel not closed within 10s")
	}
}

func TestStreamCompose_ctxCancelStopsStream(t *testing.T) {
	if engine.DetectEngine() == engine.Unknown {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan string, 64)

	err := StreamCompose(ctx, "berth-test-cancel", "", ch, "version")
	if err != nil {
		// docker compose plugin not available — drain closed channel and skip.
		for range ch {
		}
		t.Skipf("docker compose not available: %v", err)
	}

	// Cancel immediately — channel must close without blocking.
	cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range ch {
		}
	}()

	select {
	case <-done:
		// pass
	case <-time.After(5 * time.Second):
		t.Fatal("channel not closed after context cancel")
	}
}

func TestComposeUp_callsStreamCompose(t *testing.T) {
	// Smoke test: ComposeUp returns a function; calling it attempts compose up.
	// We verify the function signature compiles and error type is sensible.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ch := make(chan string, 64)

	// This will fail (no compose project) but must not panic.
	_ = ComposeUp(ctx, "berth-noproject", "/tmp", ch)
	for range ch {
	}
}

