package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Dir returns the xpoz data directory.
// On Windows: %APPDATA%\xpoz
// On Unix:    ~/.xpoz
func Dir() string {
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "xpoz")
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".xpoz")
}

// Path returns the full path to config.toml.
func Path() string {
	return filepath.Join(Dir(), "config.toml")
}

// DBPath returns the full path to tunnels.db.
func DBPath() string {
	return filepath.Join(Dir(), "tunnels.db")
}

// LogDir returns the directory for per-tunnel log files.
func LogDir() string {
	return filepath.Join(Dir(), "logs")
}
