package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	flagName  string
	flagReset bool
	flagBg    bool
)

func init() {
	rootCmd.Flags().StringVar(&flagName, "name", "", "assign a persistent name to the tunnel")
	rootCmd.Flags().BoolVar(&flagReset, "reset", false, "regenerate subdomain even if name already exists")
	rootCmd.Flags().BoolVar(&flagBg, "bg", false, "run in background (daemon mode)")
}

func runExpose(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %q — must be a number between 1 and 65535", args[0])
	}

	// TODO: implement expose logic in Phase 1
	fmt.Fprintf(cmd.OutOrStdout(), "exposing port %d...\n", port)
	return nil
}
