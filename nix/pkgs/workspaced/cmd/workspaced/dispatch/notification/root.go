package notification

import (
	"workspaced/pkg/notification"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	var title string
	var message string
	var icon string
	var urgency string
	var progress float64

	cmd := &cobra.Command{
		Use:   "notification",
		Short: "Send a desktop notification",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			n := &notification.Notification{
				Title:    title,
				Message:  message,
				Icon:     icon,
				Urgency:  urgency,
				Progress: progress,
			}
			return notification.Notify(ctx, n)
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "Workspaced", "Notification title")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Notification message")
	cmd.Flags().StringVarP(&icon, "icon", "i", "", "Notification icon")
	cmd.Flags().StringVarP(&urgency, "urgency", "u", "normal", "Notification urgency (low, normal, critical)")
	cmd.Flags().Float64VarP(&progress, "progress", "p", 0, "Notification progress (0.0-1.0)")

	return cmd
}
