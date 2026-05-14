package cli

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved tunnels",
	Args:  cobra.NoArgs,
	RunE:  runList,
}

func runList(cmd *cobra.Command, _ []string) error {
	s, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer s.Close()

	tunnels, err := s.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing tunnels: %w", err)
	}

	out := cmd.OutOrStdout()
	if len(tunnels) == 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, dimStyle.Render("  No saved tunnels."))
		fmt.Fprintln(out, dimStyle.Render("  Run 'xpoz [port] --name <name>' to create a persistent tunnel."))
		fmt.Fprintln(out)
		return nil
	}

	printTunnelTable(out, tunnels)
	return nil
}

func printTunnelTable(w io.Writer, tunnels []*store.Tunnel) {
	fmt.Fprintln(w)
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, boldStyle.Render("  NAME\tURL\tPORT\tSTATUS\tLAST USED"))
	fmt.Fprintln(tw, dimStyle.Render("  ────\t───\t────\t──────\t─────────"))
	for _, t := range tunnels {
		status := "idle"
		if t.Active {
			status = successStyle.Render("active")
		}
		fmt.Fprintf(tw, "  %s\t%s\t%d\t%s\t%s\n",
			t.Name, t.FullURL, t.LocalPort, status, dimStyle.Render(formatAge(t.LastUsed)))
	}
	tw.Flush()
	fmt.Fprintln(w)
}

// formatAge returns a human-friendly relative time string.
func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("Jan 2 2006")
	}
}
