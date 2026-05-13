package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the xpoz system service",
}

var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Register xpoz as a system service (auto-start on boot)",
	Args:  cobra.NoArgs,
	RunE:  runServiceInstall,
}

var serviceUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the xpoz system service",
	Args:  cobra.NoArgs,
	RunE:  runServiceUninstall,
}

func init() {
	serviceCmd.AddCommand(serviceInstallCmd, serviceUninstallCmd)
}

func runServiceInstall(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "service install not yet implemented")
	return nil
}

func runServiceUninstall(cmd *cobra.Command, args []string) error {
	// TODO: implement in Phase 2
	fmt.Fprintln(cmd.OutOrStdout(), "service uninstall not yet implemented")
	return nil
}
