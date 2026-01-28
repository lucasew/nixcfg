package demo

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "debug",
			Short: "Debug flag passing",
			Run: func(cmd *cobra.Command, args []string) {
				testFlag, _ := cmd.Flags().GetString("test")
				fmt.Printf("test flag value: %s\n", testFlag)
				fmt.Printf("args: %v\n", args)
			},
		}
		cmd.Flags().String("test", "default", "a test flag")
		parent.AddCommand(cmd)
	})
}
