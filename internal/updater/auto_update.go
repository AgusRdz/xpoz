package updater

// AutoUpdateConfig stores the user preference for background update checks.
type AutoUpdateConfig struct {
	Enabled bool `toml:"enabled"`
}

// GetAutoUpdateConfig reads the stored preference. Defaults to enabled.
func GetAutoUpdateConfig() AutoUpdateConfig {
	// TODO: implement in Phase 2
	return AutoUpdateConfig{Enabled: true}
}

// SetAutoUpdate persists the auto-update preference to disk.
func SetAutoUpdate(enabled bool) error {
	// TODO: implement in Phase 2
	return nil
}
