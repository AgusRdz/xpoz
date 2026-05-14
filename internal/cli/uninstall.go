package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AgusRdz/xpoz/internal/cleanup"
	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/spf13/cobra"
)

var flagKeepData bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove xpoz and all its data",
	Long: `Remove the xpoz binary, configuration, and tunnel database.

Use --keep-data to preserve the tunnel history database.`,
	Args: cobra.NoArgs,
	RunE: runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&flagKeepData, "keep-data", false, "preserve the tunnel history database")
}

func runUninstall(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	in := bufio.NewReader(cmd.InOrStdin())

	fmt.Fprintln(out)
	fmt.Fprintln(out, boldStyle.Render("xpoz uninstall"))
	fmt.Fprintln(out)

	fmt.Fprintln(out, "The following will be removed:")
	fmt.Fprintln(out, dimStyle.Render("  "+config.Path()))
	if !flagKeepData {
		fmt.Fprintln(out, dimStyle.Render("  "+config.DBPath()))
	}
	fmt.Fprintln(out)

	if !confirmPrompt(in, out, "Continue?") {
		fmt.Fprintln(out, "Cancelled.")
		return nil
	}

	binPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("uninstall: resolving binary path: %w", err)
	}

	if err := cleanup.Run(cleanup.Options{KeepData: flagKeepData}); err != nil {
		return err
	}

	if err := os.Remove(binPath); err != nil && !os.IsNotExist(err) {
		printWarn(out, fmt.Sprintf("Could not remove binary at %s: %v", binPath, err))
		printWarn(out, "You may need to remove it manually.")
	}

	printSuccess(out, "xpoz uninstalled.")
	return nil
}
