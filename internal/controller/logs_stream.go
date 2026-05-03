package controller

import (
	"context"
	"sync"
)

// LogEntry carries a single log line with its source container name.
type LogEntry struct {
	ContainerName string
	Line          string
}

// StreamMultiContainerLogs fans out one goroutine per container, writing LogEntry
// values to ch. ch is closed when all goroutines finish or ctx is cancelled.
func StreamMultiContainerLogs(ctx context.Context, containers []Container, ch chan<- LogEntry) {
	defer close(ch)

	var wg sync.WaitGroup
	for _, c := range containers {
		wg.Add(1)
		go func(c Container) {
			defer wg.Done()
			lineCh := make(chan string, 100)
			go StreamContainerLogs(ctx, c.ID, lineCh)
			for line := range lineCh {
				select {
				case <-ctx.Done():
					return
				case ch <- LogEntry{ContainerName: c.Names, Line: line}:
				}
			}
		}(c)
	}
	wg.Wait()
}
