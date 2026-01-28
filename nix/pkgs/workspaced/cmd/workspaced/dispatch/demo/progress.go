package demo

import (
	"fmt"
	"time"
	"workspaced/pkg/drivers/notification"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "progress",
			Short: "Demo progress notification",
			Run: func(cmd *cobra.Command, args []string) {
				ctx := cmd.Context()
				n := &notification.Notification{
					Title: "Demo de Progresso",
					Icon:  "utilities-terminal",
				}
				for i := 1; i <= 10; i++ {
					percent := i * 10
					n.Message = fmt.Sprintf("Passo %d de 10...", i)
					n.Progress = percent
					n.Notify(ctx)
					time.Sleep(time.Second)
				}
				n.Message = "Demo concluída!"
				n.Progress = 0
				n.Notify(ctx)
			},
		})
	})
}
