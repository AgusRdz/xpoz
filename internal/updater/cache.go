package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AgusRdz/xpoz/internal/config"
)

func cachePath() string {
	return filepath.Join(config.Dir(), "update_check")
}

// writeCache persists the latest known version tag and a timestamp.
func writeCache(tag string) error {
	line := fmt.Sprintf("%s %d\n", tag, time.Now().Unix())
	return os.WriteFile(cachePath(), []byte(line), 0600)
}

// readCache returns the cached tag and when it was written.
func readCache() (tag string, at time.Time, err error) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return "", time.Time{}, err
	}
	parts := strings.Fields(strings.TrimSpace(string(data)))
	if len(parts) != 2 {
		return "", time.Time{}, fmt.Errorf("invalid cache format")
	}
	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid timestamp: %w", err)
	}
	return parts[0], time.Unix(ts, 0), nil
}
