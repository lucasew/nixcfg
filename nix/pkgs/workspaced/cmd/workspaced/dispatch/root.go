package dispatch

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch workspace commands",
}
