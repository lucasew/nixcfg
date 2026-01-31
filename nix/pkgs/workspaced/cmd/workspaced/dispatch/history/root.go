package history

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"workspaced/pkg/db"
	"workspaced/pkg/types"

	"github.com/gorilla/websocket"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "History management",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "search [query]",
		Short: "Search history using fuzzy finder",
		RunE: func(c *cobra.Command, args []string) error {
			database, ok := c.Context().Value("db").(*db.DB)
			if !ok {
				var err error
				database, err = db.Open()
				if err != nil {
					return err
				}
				defer database.Close()
			}

			events, err := database.SearchHistory(c.Context(), "", 5000)
			if err != nil {
				return fmt.Errorf("failed to fetch history: %w", err)
			}

			if len(events) == 0 {
				return fmt.Errorf("no history found")
			}

			// 2. Run fuzzy finder
			options := []fuzzyfinder.Option{
				fuzzyfinder.WithPreviewWindow(func(i int, width int, height int) string {
					if i == -1 {
						return ""
					}
					e := events[i]
					t := time.Unix(e.Timestamp, 0).Format("2006-01-02 15:04:05")
					return fmt.Sprintf("Time:     %s\nExitCode: %d\nCwd:      %s\nDuration: %dms\n\nCommand:\n%s",
						t, e.ExitCode, e.Cwd, e.Duration, e.Command)
				}),
			}

			if len(args) > 0 {
				options = append(options, fuzzyfinder.WithQuery(strings.Join(args, " ")))
			}

			idx, err := fuzzyfinder.Find(
				events,
				func(i int) string {
					return events[i].Command
				},
				options...,
			)

			if err != nil {
				if err == fuzzyfinder.ErrAbort {
					return nil
				}
				return fmt.Errorf("fuzzy finder failed: %w", err)
			}

			fmt.Print(events[idx].Command)
			return nil
		},
	})

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List history entries (internal use)",
		RunE: func(c *cobra.Command, args []string) error {
			limit, _ := c.Flags().GetInt32("limit")
			asJSON, _ := c.Flags().GetBool("json")

			database, ok := c.Context().Value("db").(*db.DB)
			if !ok {
				var err error
				database, err = db.Open()
				if err != nil {
					return err
				}
				defer database.Close()
			}

			events, err := database.SearchHistory(c.Context(), "", int(limit))

			if err != nil {
				return err
			}

			if asJSON {
				return json.NewEncoder(c.OutOrStdout()).Encode(events)
			}

			for _, e := range events {
				t := time.Unix(e.Timestamp, 0).Format("2006-01-02 15:04:05")
				fmt.Fprintf(c.OutOrStdout(), "%s\t%s\n", t, e.Command)
			}

			return nil
		},
	}
	listCmd.Flags().Int32("limit", 5000, "Limit number of entries")
	listCmd.Flags().Bool("json", false, "Output as JSON")
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
