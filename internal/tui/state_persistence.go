package tui

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

type persistedState struct {
	CollapsedGroups map[string]bool `json:"collapsedGroups"`
}

func stateFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "berth", "state.json"), nil
}

func loadState() persistedState {
	path, err := stateFilePath()
	if err != nil {
		return persistedState{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return persistedState{}
	}
	var s persistedState
	if err := json.Unmarshal(data, &s); err != nil {
		slog.Debug("loadState: parse error", "err", err)
		return persistedState{}
	}
	return s
}

func loadedCollapsedGroups() map[string]bool {
	if m := loadState().CollapsedGroups; m != nil {
		return m
	}
	return make(map[string]bool)
}

func saveState(s persistedState) {
	path, err := stateFilePath()
	if err != nil {
		slog.Debug("saveState: cache dir error", "err", err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		slog.Debug("saveState: mkdir error", "err", err)
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		slog.Debug("saveState: marshal error", "err", err)
		return
	}
	if err := os.WriteFile(path, data, 0o640); err != nil {
		slog.Debug("saveState: write error", "err", err)
	}
}
