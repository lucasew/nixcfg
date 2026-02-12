package doctor

import (
	"fmt"
	"os"
	"text/tabwriter"
	"workspaced/pkg/driver"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "doctor",
	Short: "Check status of all registered drivers",
	Run: func(cmd *cobra.Command, args []string) {
		report := driver.Doctor(cmd.Context())

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "INTERFACE\tDRIVER\tSTATUS\tMESSAGE")

		for _, iface := range report {
			for _, d := range iface.Drivers {
				status := "❌ Unavailable"
				msg := ""
				if d.Available {
					status = "✅ Available"
				}
				if d.Error != nil {
					msg = d.Error.Error()
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", iface.Name, d.Name, status, msg)
			}
		}
		w.Flush()
	},
}
