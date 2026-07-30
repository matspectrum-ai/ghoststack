package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type RemoteConfig struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	APIKey string `json:"api_key"`
}

const remoteDir = ".ghoststack"
const remoteFile = "remotes.json"

func remotesPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, remoteDir, remoteFile)
}

func loadRemotes() ([]RemoteConfig, error) {
	path := remotesPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []RemoteConfig{}, nil
		}
		return nil, err
	}
	var remotes []RemoteConfig
	if err := json.Unmarshal(data, &remotes); err != nil {
		return nil, err
	}
	return remotes, nil
}

func saveRemotes(remotes []RemoteConfig) error {
	path := remotesPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(remotes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func findRemote(name string) (*RemoteConfig, error) {
	remotes, err := loadRemotes()
	if err != nil {
		return nil, err
	}
	for _, r := range remotes {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("remote %q not found", name)
}

func newRemoteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage remote GhostStack instances",
	}

	cmd.AddCommand(newRemoteAddCommand())
	cmd.AddCommand(newRemoteListCommand())
	cmd.AddCommand(newRemoteRemoveCommand())
	cmd.AddCommand(newRemoteExecCommand())
	cmd.AddCommand(newRemoteLogsCommand())

	return cmd
}

func newRemoteAddCommand() *cobra.Command {
	var apiKey string
	cmd := &cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add a remote GhostStack instance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiKey == "" {
				return fmt.Errorf("--api-key is required")
			}

			remotes, err := loadRemotes()
			if err != nil {
				return err
			}

			for _, r := range remotes {
				if r.Name == args[0] {
					return fmt.Errorf("remote %q already exists", args[0])
				}
			}

			remotes = append(remotes, RemoteConfig{
				Name:   args[0],
				URL:    strings.TrimRight(args[1], "/"),
				APIKey: apiKey,
			})

			if err := saveRemotes(remotes); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Remote %q added (%s)\n", args[0], args[1])
			return nil
		},
	}
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	return cmd
}

func newRemoteListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured remotes",
		RunE: func(cmd *cobra.Command, args []string) error {
			remotes, err := loadRemotes()
			if err != nil {
				return err
			}

			if len(remotes) == 0 {
				fmt.Fprintln(os.Stdout, "No remotes configured.")
				return nil
			}

			for _, r := range remotes {
				fmt.Fprintf(os.Stdout, "  %-20s %s\n", r.Name, r.URL)
			}
			return nil
		},
	}
}

func newRemoteRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove a remote",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			remotes, err := loadRemotes()
			if err != nil {
				return err
			}

			found := false
			var updated []RemoteConfig
			for _, r := range remotes {
				if r.Name == args[0] {
					found = true
					continue
				}
				updated = append(updated, r)
			}

			if !found {
				return fmt.Errorf("remote %q not found", args[0])
			}

			if err := saveRemotes(updated); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Remote %q removed.\n", args[0])
			return nil
		},
	}
}

func newRemoteExecCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exec <remote> <command>",
		Short: "Execute a command on a remote instance (status, logs, start, stop, restart)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			remote, err := findRemote(args[0])
			if err != nil {
				return err
			}

			client := &http.Client{Timeout: 15 * time.Second}

			var url string
			switch args[1] {
			case "status":
				url = remote.URL + "/api/status"
			case "logs":
				url = remote.URL + "/api/logs"
			default:
				return fmt.Errorf("unknown command: %s (use: status, logs)", args[1])
			}

			req, err := http.NewRequestWithContext(cmd.Context(), "GET", url, nil)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Bearer "+remote.APIKey)

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("request: %w", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("remote returned %d: %s", resp.StatusCode, string(body))
			}

			var pretty bytes.Buffer
			if err := json.Indent(&pretty, body, "", "  "); err != nil {
				fmt.Fprintln(os.Stdout, string(body))
			} else {
				fmt.Fprintln(os.Stdout, pretty.String())
			}

			return nil
		},
	}
}

func newRemoteLogsCommand() *cobra.Command {
	var follow bool
	cmd := &cobra.Command{
		Use:   "logs <remote>",
		Short: "Fetch logs from a remote instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newRemoteExecCommand().RunE(cmd, append(args, "logs"))
		},
	}
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output")
	return cmd
}
