package svc

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"workspaced/pkg/drivers/screen"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "screencaps",
			Short: "Monitor CapsLock and toggle screen DPMS",
			Run: func(cmd *cobra.Command, args []string) {
				monitorCapsLock(cmd.Context())
			},
		})
	})
}

func monitorCapsLock(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	matches, _ := filepath.Glob("/sys/class/leds/*capslock/brightness")
	if len(matches) == 0 {
		slog.Warn("no capslock leds found")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			capsActive := false
			for _, m := range matches {
				data, err := os.ReadFile(m)
				if err == nil && strings.TrimSpace(string(data)) == "1" {
					capsActive = true
					break
				}
			}

			screenActive, err := screen.IsDPMSOn(ctx)
			if err != nil {
				slog.Error("on checking if screen is active", "error", err)
			}
			if !capsActive != screenActive {
				slog.Info("toggling screen", "active", !capsActive)
				screen.SetDPMS(ctx, !capsActive)
			}
		}
	}
}
