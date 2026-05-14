package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/AgusRdz/xpoz/internal/config"
)

var configureCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or edit xpoz configuration",
	Long: `Show the current xpoz configuration or open it in $EDITOR.

Subcommands:
  xpoz config          show config file path and contents
  xpoz config edit     open config in $EDITOR`,
	Args: cobra.NoArgs,
	RunE: runConfigShow,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config in $EDITOR",
	Args:  cobra.NoArgs,
	RunE:  runConfigEdit,
}

func init() {
	configureCmd.AddCommand(configEditCmd)
}

func runConfigShow(cmd *cobra.Command, _ []string) error {
	p := config.Path()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no config file found at %s — run 'xpoz setup' first", p)
		}
		return fmt.Errorf("config: reading file: %w", err)
	}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, dimStyle.Render(p))
	fmt.Fprintln(out)
	fmt.Fprint(out, string(data))
	return nil
}

func runConfigEdit(cmd *cobra.Command, _ []string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vi"
		}
	}
	c := exec.CommandContext(cmd.Context(), editor, config.Path())
	c.Stdin = cmd.InOrStdin()
	c.Stdout = cmd.OutOrStdout()
	c.Stderr = cmd.ErrOrStderr()
	return c.Run()
}
