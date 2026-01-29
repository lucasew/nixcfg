package registry

import "github.com/spf13/cobra"

// RegisterFunc is a function type that modifies a cobra.Command.
// Implementations typically add subcommands or flags to the provided parent command.
type RegisterFunc func(*cobra.Command)

// CommandRegistry is a helper to aggregate command builders.
// It allows different packages to register their subcommands independently,
// facilitating a modular CLI structure without cyclic dependencies.
type CommandRegistry struct {
	builders []RegisterFunc
}

// Register adds a new builder function to the registry.
// The provided function 'f' will be executed later when GetCommand is invoked.
func (r *CommandRegistry) Register(f RegisterFunc) {
	r.builders = append(r.builders, f)
}

// GetCommand applies all registered builder functions to the base command.
// It iterates through all registered functions, allowing them to attach their
// respective subcommands to the 'base' command.
func (r *CommandRegistry) GetCommand(base *cobra.Command) *cobra.Command {
	for _, build := range r.builders {
		build(base)
	}
	return base
}
