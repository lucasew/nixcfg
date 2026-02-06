package svc

import (
	"context"
	"log/slog"
	"strings"
	"time"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "renice-hungry",
			Short: "Lowers the priority of the most cpu hungry process periodically",
			Run: func(cmd *cobra.Command, args []string) {
				ctx := cmd.Context()
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()

				slog.Info("renice-hungry started")

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						pid, cmdline, err := getHungryPID(ctx)
						if err != nil {
							slog.Error("failed to get hungry PID", "error", err)
							continue
						}
						if pid == "" {
							continue
						}

						slog.Info("renicing process", "pid", pid, "cmd", cmdline)
						_ = exec.RunCmd(ctx, "renice", "7", pid).Run()
					}
				}
			},
		})
	})
}

func getHungryPID(ctx context.Context) (string, string, error) {
	// ps -eo pid,args --sort=-%cpu | head -n2 | tail -n 1
	out, err := exec.RunCmd(ctx, "ps", "-eo", "pid,args", "--sort=-%cpu").Output()
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return "", "", nil
	}

	// First line is header, second line is top process
	line := strings.TrimSpace(lines[1])
	parts := strings.SplitN(line, " ", 2)
	pid := parts[0]
	cmdline := ""
	if len(parts) > 1 {
		cmdline = strings.TrimSpace(parts[1])
	}
	return pid, cmdline, nil
}
