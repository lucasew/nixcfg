package svc

import (
	"fmt"
	"log/slog"
	"time"
	"workspaced/pkg/drivers/battery"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "osmardetector",
			Short: "Annoying beep each second if laptop stops charging",
			Run: func(cmd *cobra.Command, args []string) {
				ctx := cmd.Context()
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()

				slog.Info("osmardetector started")

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						status, err := battery.GetStatus(ctx)
						if err != nil {
							slog.Error("failed to get battery status", "error", err)
							continue
						}
						if status == battery.Discharging {
							fmt.Print("\aAi!")
						}
					}
				}
			},
		})
	})
}
