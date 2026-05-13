package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and apply the latest xpoz release",
	Args:  cobra.NoArgs,
	RunE:  runUpdate,
}

var autoUpdateCmd = &cobra.Command{
	Use:       "auto-update [on|off]",
	Short:     "Enable or disable automatic background updates",
	ValidArgs: []string{"on", "off"},
	Args:      cobra.MaximumNArgs(1),
	RunE:      runAutoUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "update not yet implemented")
	return nil
}

func runAutoUpdate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// TODO: show current status
		fmt.Fprintln(cmd.OutOrStdout(), "auto-update status: not yet implemented")
		return nil
	}
	switch args[0] {
	case "on", "off":
		// TODO: implement
		fmt.Fprintf(cmd.OutOrStdout(), "auto-update %s: not yet implemented\n", args[0])
	default:
		return fmt.Errorf("unknown argument %q — use 'on' or 'off'", args[0])
	}
	return nil
}
