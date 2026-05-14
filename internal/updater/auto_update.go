package updater

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AgusRdz/xpoz/internal/config"
)

// AutoUpdateConfig stores the user preference for background update checks.
type AutoUpdateConfig struct {
	Enabled bool
}

func autoUpdatePath() string {
	return filepath.Join(config.Dir(), "auto_update")
}

// GetAutoUpdateConfig reads the stored preference. Defaults to enabled.
func GetAutoUpdateConfig() AutoUpdateConfig {
	data, err := os.ReadFile(autoUpdatePath())
	if err != nil {
		return AutoUpdateConfig{Enabled: true}
	}
	return AutoUpdateConfig{Enabled: strings.TrimSpace(string(data)) != "false"}
}

// SetAutoUpdate persists the auto-update preference to disk.
func SetAutoUpdate(enabled bool) error {
	val := "true"
	if !enabled {
		val = "false"
	}
	return os.WriteFile(autoUpdatePath(), []byte(val+"\n"), 0600)
}
