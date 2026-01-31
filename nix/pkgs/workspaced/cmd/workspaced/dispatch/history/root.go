package history

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"workspaced/pkg/db"
	"workspaced/pkg/types"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "History management",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "search [query]",
		Short: "Search history using fzf",
		RunE: func(c *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = strings.Join(args, " ")
			}

			// We need to get history from the daemon
			// Since we want to pipe to fzf, we can't let the normal dispatch mechanism handle it
			// unless we make 'list' a command that outputs what we need.

			// 1. Get history from daemon via 'history list'
			listCmd := exec.Command("workspaced", "dispatch", "history", "list", "--limit", "5000")
			out, err := listCmd.Output()
			if err != nil {
				return fmt.Errorf("failed to list history: %w", err)
			}

			// 2. Run fzf
			fzfCmd := exec.Command("fzf",
				"--delimiter", "\t",
				"--with-nth", "2..",
				"--query", query,
				"--layout", "reverse",
				"--height", "40%",
				"--header", "Workspaced History",
				"--preview", "echo -e \"Time: {1}\nCommand: {2..}\"",
				"--preview-window", "bottom:3:wrap",
			)
			fzfCmd.Stdin = strings.NewReader(string(out))
			fzfCmd.Stderr = os.Stderr

			selection, err := fzfCmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
					return nil
				}
				return fmt.Errorf("fzf failed: %w", err)
			}

			selectedStr := strings.TrimSpace(string(selection))
			if selectedStr != "" {
				parts := strings.Split(selectedStr, "\t")
				if len(parts) >= 2 {
					fmt.Print(strings.Join(parts[1:], "\t"))
				}
			}

			return nil
		},
	})

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List history entries (internal use)",
		RunE: func(c *cobra.Command, args []string) error {
			limit, _ := c.Flags().GetInt32("limit")

			database, ok := c.Context().Value("db").(*db.DB)
			if !ok {
				// If not on daemon, this will be dispatched to daemon
				return fmt.Errorf("this command must run on the daemon")
			}

			events, err := database.SearchHistory(c.Context(), "", int(limit))
			if err != nil {
				return err
			}

			for _, e := range events {
				t := time.Unix(e.Timestamp, 0).Format("2006-01-02 15:04:05")
				fmt.Fprintf(c.OutOrStdout(), "%s\t%s\n", t, e.Command)
			}

			return nil
		},
	}
	listCmd.Flags().Int32("limit", 5000, "Limit number of entries")
	cmd.AddCommand(listCmd)

	recordCmd := &cobra.Command{
		Use:   "record",
		Short: "Record a command in history",
		RunE: func(c *cobra.Command, args []string) error {
			var event types.HistoryEvent

			// Try reading from stdin if no command flag is provided
			command, _ := c.Flags().GetString("command")
			if command == "" {
				if err := json.NewDecoder(os.Stdin).Decode(&event); err != nil {
					return err
				}
			} else {
				event.Command = command
				event.Cwd, _ = c.Flags().GetString("cwd")
				event.ExitCode, _ = c.Flags().GetInt("exit-code")
				event.Timestamp, _ = c.Flags().GetInt64("timestamp")
				event.Duration, _ = c.Flags().GetInt64("duration")
			}

			if event.Timestamp == 0 {
				event.Timestamp = time.Now().Unix()
			}
			if event.Cwd == "" {
				event.Cwd, _ = os.Getwd()
			}

			if database, ok := c.Context().Value("db").(*db.DB); ok {
				return database.RecordHistory(c.Context(), event)
			}

			return sendHistoryEvent(event)
		},
	}
	recordCmd.Flags().String("command", "", "Command string")
	recordCmd.Flags().String("cwd", "", "Current working directory")
	recordCmd.Flags().Int("exit-code", 0, "Exit code")
	recordCmd.Flags().Int64("timestamp", 0, "Timestamp")
	recordCmd.Flags().Int64("duration", 0, "Duration in ms")
	cmd.AddCommand(recordCmd)

	return cmd
}

func getSocketPath() string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return filepath.Join(runtimeDir, "workspaced.sock")
}

func sendHistoryEvent(event types.HistoryEvent) error {
	socketPath := getSocketPath()
	dialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, 2*time.Second)
		},
	}

	conn, _, err := dialer.Dial("ws://localhost/ws", nil)
	if err != nil {
		return nil // Daemon not running, ignore
	}
	defer conn.Close()

	payload, _ := json.Marshal(event)
	packet := types.StreamPacket{
		Type:    "history_event",
		Payload: payload,
	}

	return conn.WriteJSON(packet)
}
