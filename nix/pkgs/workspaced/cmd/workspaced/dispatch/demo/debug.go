package demo

import (
	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "debug",
			Short: "Debug flag passing",
			Run: func(cmd *cobra.Command, args []string) {
				testFlag, _ := cmd.Flags().GetString("test")
				cmd.Printf("test flag value: %s\n", testFlag)
				cmd.Printf("args: %v\n", args)
			},
		}
		cmd.Flags().String("test", "default", "a test flag")
		parent.AddCommand(cmd)
	})
}
