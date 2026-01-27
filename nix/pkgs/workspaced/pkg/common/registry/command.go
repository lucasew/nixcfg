package registry

import "github.com/spf13/cobra"

type RegisterFunc func(*cobra.Command)

type CommandRegistry struct {
	builders []RegisterFunc
}

func (r *CommandRegistry) Register(f RegisterFunc) {
	r.builders = append(r.builders, f)
}

func (r *CommandRegistry) GetCommand(base *cobra.Command) *cobra.Command {
	for _, build := range r.builders {
		build(base)
	}
	return base
}
