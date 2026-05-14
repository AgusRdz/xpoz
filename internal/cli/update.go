package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AgusRdz/xpoz/internal/updater"
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

func runUpdate(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	printInfo(out, "Checking for updates...")
	if err := updater.Update(cmd.Context(), rootCmd.Version); err != nil {
		// "already on latest" is not a fatal error — print it as info
		if rootCmd.Version != "dev" {
			printInfo(out, err.Error())
			return nil
		}
		return err
	}
	printSuccess(out, "xpoz updated — restart to use the new version.")
	return nil
}

func runAutoUpdate(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()
	if len(args) == 0 {
		cfg := updater.GetAutoUpdateConfig()
		status := "on"
		if !cfg.Enabled {
			status = "off"
		}
		fmt.Fprintf(out, "auto-update: %s\n", status)
		return nil
	}
	switch args[0] {
	case "on":
		if err := updater.SetAutoUpdate(true); err != nil {
			return fmt.Errorf("enabling auto-update: %w", err)
		}
		printSuccess(out, "Auto-update enabled.")
	case "off":
		if err := updater.SetAutoUpdate(false); err != nil {
			return fmt.Errorf("disabling auto-update: %w", err)
		}
		printSuccess(out, "Auto-update disabled.")
	default:
		return fmt.Errorf("unknown argument %q — use 'on' or 'off'", args[0])
	}
	return nil
}
