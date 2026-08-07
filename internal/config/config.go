//go:build windows

// Package config persists settings to config.json next to the exe, keeping the
// app portable.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// exePathOnce / exePath cache the config file location: os.Executable is a
// syscall and would otherwise run on every Load/Save.
var (
	exePathOnce sync.Once
	exePath     string
)

// PersistConfig mirrors the panel controls.
// Note: run-at-startup is a system-level setting (registry Run key), not part of
// the config file; see the registry package.
type PersistConfig struct {
	DelayMS        int    `json:"delay_ms"`
	AutoStartGuard bool   `json:"auto_start_guard"` // auto-guard after login
	MinToTray      bool   `json:"min_to_tray"`      // minimize to tray
	CloseToTray    bool   `json:"close_to_tray"`    // close window minimizes to tray
	Lang           string `json:"lang"`             // UI language code; empty = auto (zh/en/ja/ko/fr/de/es/ru/pt/it)
}

// Default returns the default configuration.
func Default() PersistConfig {
	return PersistConfig{
		DelayMS:        3000,
		AutoStartGuard: false,
		MinToTray:      true,
		CloseToTray:    true,
	}
}

// Path returns the config file path next to the exe (cached after the first
// call).
func Path() string {
	exePathOnce.Do(func() {
		exe, err := os.Executable()
		if err != nil {
			exePath = "config.json"
			return
		}
		exePath = filepath.Join(filepath.Dir(exe), "config.json")
	})
	return exePath
}

// Load reads the config; returns defaults when the file is missing or invalid.
func Load() PersistConfig {
	c := Default()
	b, err := os.ReadFile(Path())
	if err != nil {
		return c
	}
	if json.Unmarshal(b, &c) != nil {
		return Default()
	}
	if c.DelayMS < 100 {
		c.DelayMS = 100
	}
	return c
}

// Save writes the config atomically (temp file + rename), so a crash mid-write
// can never leave a truncated/corrupt config.json; returns an error for the
// caller to report.
func Save(c PersistConfig) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	p := Path()
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	// os.Rename on Windows fails when the destination exists, so remove it
	// first; the temp file keeps the data safe until the swap.
	if err := os.Rename(tmp, p); err != nil {
		if rmErr := os.Remove(p); rmErr != nil {
			os.Remove(tmp)
			return rmErr
		}
		if err := os.Rename(tmp, p); err != nil {
			os.Remove(tmp)
			return err
		}
	}
	return nil
}
