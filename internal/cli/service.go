package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/AgusRdz/xpoz/internal/config"
	"github.com/AgusRdz/xpoz/internal/proxy"
	"github.com/AgusRdz/xpoz/internal/store"
	svc "github.com/kardianos/service"
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

var serviceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the xpoz system service",
	Args:  cobra.NoArgs,
	RunE:  runServiceStart,
}

var serviceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the xpoz system service",
	Args:  cobra.NoArgs,
	RunE:  runServiceStop,
}

// serviceRunCmd is called by the OS service manager — not for direct user use.
var serviceRunCmd = &cobra.Command{
	Use:    "run",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE:   runServiceRun,
}

func init() {
	serviceCmd.AddCommand(
		serviceInstallCmd,
		serviceUninstallCmd,
		serviceStartCmd,
		serviceStopCmd,
		serviceRunCmd,
	)
}

// xpozSvcConfig returns the kardianos/service configuration for xpoz.
func xpozSvcConfig() *svc.Config {
	return &svc.Config{
		Name:        "xpoz",
		DisplayName: "xpoz tunnel service",
		Description: "Exposes local ports via your own domain and Cloudflare Tunnel.",
		Arguments:   []string{"service", "run"},
	}
}

func newService() (svc.Service, error) {
	prog := &svcProgram{}
	s, err := svc.New(prog, xpozSvcConfig())
	if err != nil {
		return nil, fmt.Errorf("creating service handle: %w", err)
	}
	prog.svc = s
	return s, nil
}

func runServiceInstall(cmd *cobra.Command, _ []string) error {
	s, err := newService()
	if err != nil {
		return err
	}
	if err := s.Install(); err != nil {
		return fmt.Errorf("installing service: %w", err)
	}
	if err := s.Start(); err != nil {
		return fmt.Errorf("starting service after install: %w", err)
	}
	printSuccess(cmd.OutOrStdout(), "xpoz service installed and started.")
	return nil
}

func runServiceUninstall(cmd *cobra.Command, _ []string) error {
	s, err := newService()
	if err != nil {
		return err
	}
	_ = s.Stop() // best-effort; may already be stopped
	if err := s.Uninstall(); err != nil {
		return fmt.Errorf("uninstalling service: %w", err)
	}
	printSuccess(cmd.OutOrStdout(), "xpoz service uninstalled.")
	return nil
}

func runServiceStart(cmd *cobra.Command, _ []string) error {
	s, err := newService()
	if err != nil {
		return err
	}
	if err := s.Start(); err != nil {
		return fmt.Errorf("starting service: %w", err)
	}
	printSuccess(cmd.OutOrStdout(), "xpoz service started.")
	return nil
}

func runServiceStop(cmd *cobra.Command, _ []string) error {
	s, err := newService()
	if err != nil {
		return err
	}
	if err := s.Stop(); err != nil {
		return fmt.Errorf("stopping service: %w", err)
	}
	printSuccess(cmd.OutOrStdout(), "xpoz service stopped.")
	return nil
}

func runServiceRun(_ *cobra.Command, _ []string) error {
	s, err := newService()
	if err != nil {
		return err
	}
	return s.Run()
}

// svcProgram implements svc.Interface and holds live service state.
type svcProgram struct {
	svc    svc.Service
	prx    *proxy.Server
	db     store.Store
	cancel context.CancelFunc
}

func (p *svcProgram) Start(_ svc.Service) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("service start: loading config: %w", err)
	}
	db, err := store.New(config.DBPath())
	if err != nil {
		return fmt.Errorf("service start: opening store: %w", err)
	}
	p.db = db
	p.prx = proxy.New(cfg.Proxy.Port)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go p.run(ctx, cfg)
	return nil
}

func (p *svcProgram) Stop(_ svc.Service) error {
	if p.cancel != nil {
		p.cancel()
	}
	if p.prx != nil {
		_ = p.prx.Stop()
	}
	if p.db != nil {
		_ = p.db.Close()
	}
	return nil
}

func (p *svcProgram) run(ctx context.Context, cfg *config.Config) {
	if err := p.loadRoutes(ctx, cfg); err != nil {
		log.Printf("xpoz service: loading routes: %v", err)
	}
	if err := p.prx.Start(); err != nil {
		log.Printf("xpoz service: proxy: %v", err)
	}
}

func (p *svcProgram) loadRoutes(ctx context.Context, cfg *config.Config) error {
	tunnels, err := p.db.List(ctx)
	if err != nil {
		return fmt.Errorf("listing tunnels: %w", err)
	}
	for _, t := range tunnels {
		hostname := t.Subdomain + "." + cfg.Domain.Name
		p.prx.AddRoute(hostname, fmt.Sprintf("localhost:%d", t.LocalPort))
	}
	return nil
}
