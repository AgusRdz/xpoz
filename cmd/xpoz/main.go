package main

import (
	"fmt"
	"os"

	"github.com/AgusRdz/xpoz/internal/cli"
	"github.com/AgusRdz/xpoz/internal/updater"
)

// version is set at build time via -ldflags "-X main.version=v1.2.3"
var version = "dev"

func main() {
	updater.ApplyPending(version)
	updater.NotifyIfAvailable(version)

	if err := cli.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		os.Exit(1)
	}
}
