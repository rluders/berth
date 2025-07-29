// Package main is the entry point for the Berth TUI application.
package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rluders/berth/internal/tui"
)

// main function initializes and runs the Bubble Tea program.
func main() {
	// Setup logging to a file
	logFile, err := os.OpenFile("berth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	// Recover from panics and log them
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in main defer", "panic", r)
			fmt.Printf("Alas, there's been a panic: %v\n", r)
			os.Exit(1)
		}
	}()

	var program *tea.Program
	func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered during program initialization", "panic", r)
				fmt.Printf("Alas, there's been a panic during init: %v\n", r)
				os.Exit(1)
			}
		}()
		slog.Debug("Initializing Bubble Tea program...")
		program = tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	}()

	slog.Debug("Running Bubble Tea program...")
	if _, err := program.Run(); err != nil {
		slog.Error("Program error", "error", err)
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
