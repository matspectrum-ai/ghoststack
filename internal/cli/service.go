package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const servicePath = "/etc/systemd/system/ghoststack.service"
const configDir = "/etc/ghoststack"
const dataDir = "/var/lib/ghoststack"

func newServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage systemd service",
	}

	cmd.AddCommand(newServiceInstallCommand())
	cmd.AddCommand(newServiceUninstallCommand())
	cmd.AddCommand(newServiceEnableCommand())
	cmd.AddCommand(newServiceDisableCommand())
	cmd.AddCommand(newServiceStatusCommand())

	return cmd
}

func newServiceInstallCommand() *cobra.Command {
	var user string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install ghoststack systemd service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Geteuid() != 0 {
				return fmt.Errorf("service install requires root (try sudo ghost service install)")
			}

			binary, err := os.Executable()
			if err != nil {
				return fmt.Errorf("detect binary path: %w", err)
			}

			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("create %s: %w", configDir, err)
			}
			if err := os.MkdirAll(dataDir, 0700); err != nil {
				return fmt.Errorf("create %s: %w", dataDir, err)
			}

			tplContent, err := os.ReadFile(filepath.Join(filepath.Dir(binary), "..", "share", "ghoststack", "ghoststack.service.template"))
			if err != nil {
				// fallback: embedded default template
				tplContent = []byte(defaultServiceTemplate)
			}

			tmpl, err := template.New("service").Parse(string(tplContent))
			if err != nil {
				return fmt.Errorf("parse template: %w", err)
			}

			wd, _ := os.Getwd()
			data := map[string]string{
				"BINARY":     binary,
				"USER":       user,
				"CONFIG_DIR": configDir,
				"DATA_DIR":   dataDir,
				"WORK_DIR":   wd,
			}

			if user == "" {
				user = "root"
			}
			data["USER"] = user

			if dryRun {
				if err := tmpl.Execute(os.Stdout, data); err != nil {
					return fmt.Errorf("render: %w", err)
				}
				return nil
			}

			f, err := os.Create(servicePath)
			if err != nil {
				return fmt.Errorf("create %s: %w", servicePath, err)
			}
			defer f.Close()

			if err := tmpl.Execute(f, data); err != nil {
				return fmt.Errorf("write unit: %w", err)
			}

			if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
				return fmt.Errorf("daemon-reload: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Service installed: %s\n", servicePath)
			fmt.Fprintln(os.Stdout, "Run: systemctl enable ghoststack && systemctl start ghoststack")
			return nil
		},
	}

	cmd.Flags().StringVar(&user, "user", "root", "system user for the service")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print unit without installing")
	return cmd
}

func newServiceUninstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove ghoststack systemd service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Geteuid() != 0 {
				return fmt.Errorf("service uninstall requires root (try sudo ghost service uninstall)")
			}

			exec.Command("systemctl", "stop", "ghoststack").Run()
			exec.Command("systemctl", "disable", "ghoststack").Run()

			if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove %s: %w", servicePath, err)
			}

			if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
				return fmt.Errorf("daemon-reload: %w", err)
			}

			fmt.Fprintln(os.Stdout, "Service uninstalled")
			return nil
		},
	}
}

func newServiceEnableCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable ghoststack service (start on boot)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSystemctl("enable", "ghoststack")
		},
	}
}

func newServiceDisableCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable ghoststack service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSystemctl("disable", "ghoststack")
		},
	}
}

func newServiceStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show ghoststack service status",
		RunE: func(c *cobra.Command, args []string) error {
			sh := exec.Command("systemctl", "status", "ghoststack")
			sh.Stdout = os.Stdout
			sh.Stderr = os.Stderr
			return sh.Run()
		},
	}
}

func runSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const defaultServiceTemplate = `[Unit]
Description=GhostStack Privacy Orchestration
Documentation=https://ghoststack.dev
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart={{.BINARY}} start --config {{.CONFIG_DIR}}/config.yaml
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
RestartSec=5
User={{.USER}}
LimitNOFILE=65536
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_RAW CAP_SYS_ADMIN CAP_SYS_PTRACE
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_RAW CAP_SYS_ADMIN CAP_SYS_PTRACE
DeviceAllow=/dev/net/tun rw
NoNewPrivileges=yes
ProtectSystem=full
ProtectHome=yes
ReadWritePaths={{.DATA_DIR}} {{.CONFIG_DIR}}
PrivateTmp=yes
MemoryDenyWriteExecute=yes
RestrictRealtime=yes
RestrictAddressFamilies=AF_INET AF_INET6 AF_NETLINK AF_UNIX
SystemCallArchitectures=native

[Install]
WantedBy=multi-user.target
`
