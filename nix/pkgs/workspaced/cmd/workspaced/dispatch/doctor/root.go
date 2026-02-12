package doctor

import (
	"fmt"
	"workspaced/pkg/driver"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "doctor",
	Short: "Check status of all registered drivers",
	Run: func(cmd *cobra.Command, args []string) {
		report := driver.Doctor(cmd.Context())

		fmt.Println("🩺 SYSTEM DOCTOR")
		fmt.Println("===============")

		for _, iface := range report {
			fmt.Printf("\n📦 Interface: %s\n", iface.Name)
			for _, d := range iface.Drivers {
				if d.Available {
					fmt.Printf("   ✅ %-15s: Available\n", d.Name)
				} else {
					fmt.Printf("   ❌ %-15s: %v\n", d.Name, d.Error)
				}
			}
		}
	},
}
