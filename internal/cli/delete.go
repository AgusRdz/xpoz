package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/store"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a saved tunnel by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	s, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer s.Close()

	if err := s.Delete(cmd.Context(), name); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("no tunnel named %q — run 'xpoz list' to see saved tunnels", name)
		}
		return fmt.Errorf("deleting %q: %w", name, err)
	}

	printSuccess(cmd.OutOrStdout(), fmt.Sprintf("Tunnel %q deleted.", name))
	return nil
}
