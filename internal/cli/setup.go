package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/dns"
	"github.com/AgusRdz/xpoz/internal/tunnel"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive initial configuration",
	Long:  `Configure your domain, Cloudflare API token, and Tunnel ID.`,
	Args:  cobra.NoArgs,
	RunE:  runSetup,
}

type setupInput struct {
	domain   string
	token    string
	tunnelID string
}

func runSetup(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	in := bufio.NewReader(cmd.InOrStdin())
	ctx := cmd.Context()

	if config.IsConfigured() {
		printWarn(out, "xpoz is already configured.")
		if !confirmPrompt(in, out, "Reconfigure?") {
			return nil
		}
	}

	input, err := gatherSetupInput(in, out)
	if err != nil {
		return err
	}

	printStep(out, 1, "Validating Cloudflare credentials…")
	zoneID, err := dns.GetZoneID(ctx, input.token, input.domain)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	client, err := dns.New(input.token, zoneID)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	if err := client.ValidateToken(ctx); err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	printSuccess(out, "Token valid — zone "+zoneID)

	printStep(out, 2, "Creating wildcard DNS record…")
	if err := client.EnsureWildcardCNAME(ctx, input.domain, input.tunnelID); err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	printSuccess(out, fmt.Sprintf("*.%s → %s.cfargotunnel.com", input.domain, input.tunnelID))

	printStep(out, 3, "Updating cloudflared config…")
	if err := applyCloudflaredIngress(input.domain); err != nil {
		printWarn(out, err.Error())
		printWarn(out, "Add the ingress rule manually, then re-run 'xpoz setup' to apply it.")
	} else {
		printSuccess(out, "Wildcard ingress rule added to cloudflared config")
	}

	printStep(out, 4, "Saving xpoz config…")
	cfg := &config.Config{
		Domain:     config.DomainConfig{Name: input.domain},
		Cloudflare: config.CloudflareConfig{Token: input.token, ZoneID: zoneID, Tunnel: input.tunnelID},
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	printSuccess(out, "Config saved to "+config.Path())

	printSetupComplete(out, input.domain)
	return nil
}

func gatherSetupInput(in *bufio.Reader, out io.Writer) (setupInput, error) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, boldStyle.Render("xpoz setup — connect your domain to Cloudflare Tunnel"))
	fmt.Fprintln(out)

	domain, err := promptLine(in, out, "Domain (e.g. example.com)")
	if err != nil {
		return setupInput{}, err
	}
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") {
		return setupInput{}, fmt.Errorf("setup: %q does not look like a valid domain", domain)
	}

	token, err := promptSecret(in, out, "Cloudflare API token (DNS:Edit)")
	if err != nil {
		return setupInput{}, err
	}
	if token == "" {
		return setupInput{}, fmt.Errorf("setup: Cloudflare API token cannot be empty")
	}

	tunnelID, err := promptLine(in, out, "Cloudflare Tunnel ID")
	if err != nil {
		return setupInput{}, err
	}
	if tunnelID == "" {
		return setupInput{}, fmt.Errorf("setup: Tunnel ID cannot be empty")
	}

	return setupInput{domain: domain, token: token, tunnelID: tunnelID}, nil
}

func applyCloudflaredIngress(domain string) error {
	cfgPath := tunnel.DefaultConfigPath()
	cfg, err := tunnel.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("could not read cloudflared config at %s — run 'cloudflared tunnel create' first", cfgPath)
	}
	tunnel.EnsureWildcardIngress(cfg, domain, config.DefaultProxyPort)
	if err := tunnel.Save(cfgPath, cfg); err != nil {
		return fmt.Errorf("could not write cloudflared config at %s: %w", cfgPath, err)
	}
	return nil
}

func printSetupComplete(out io.Writer, domain string) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, boldStyle.Render("Setup complete!"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Restart cloudflared to apply the ingress rule:")
	fmt.Fprintln(out, "  "+dimStyle.Render("sudo systemctl restart cloudflared"))
	fmt.Fprintln(out, "  "+dimStyle.Render("# macOS: brew services restart cloudflared"))
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Then expose a local port:")
	fmt.Fprintln(out, "  "+urlStyle.Render(fmt.Sprintf("xpoz 3000   # → https://<subdomain>.%s", domain)))
	fmt.Fprintln(out)
}
