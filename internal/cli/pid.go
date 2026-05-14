package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AgusRdz/xpoz/internal/config"
)

func pidPath(name string) string {
	return filepath.Join(config.Dir(), "pids", name+".pid")
}

func writePID(name string, pid int) error {
	path := pidPath(name)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("pid: creating dir: %w", err)
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0600)
}

func readPID(name string) (int, error) {
	data, err := os.ReadFile(pidPath(name))
	if err != nil {
		return 0, fmt.Errorf("pid: %q: %w", name, err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, fmt.Errorf("pid: invalid value in file: %w", err)
	}
	return pid, nil
}

func removePID(name string) {
	_ = os.Remove(pidPath(name))
}

// pidAlive reports whether a PID file exists for name and the process is still running.
func pidAlive(name string) bool {
	pid, err := readPID(name)
	if err != nil {
		return false
	}
	return isRunning(pid)
}
