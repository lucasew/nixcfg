package termux

import (
	"context"
	"fmt"
	"os"
	"strings"
	"workspaced/pkg/driver"
	"workspaced/pkg/driver/terminal"
	execdriver "workspaced/pkg/driver/exec"
)

func init() {
	driver.Register[terminal.Driver](&Provider{})
}

type Provider struct{}

func (p *Provider) ID() string         { return "terminal_termux" }
func (p *Provider) Name() string       { return "Termux" }
func (p *Provider) DefaultWeight() int { return 0 }

func (p *Provider) CheckCompatibility(ctx context.Context) error {
	if os.Getenv("TERMUX_VERSION") == "" {
		return fmt.Errorf("%w: not running in Termux", driver.ErrIncompatible)
	}
	return nil
}

func (p *Provider) New(ctx context.Context) (terminal.Driver, error) {
	return &Driver{}, nil
}

type Driver struct{}

func (d *Driver) Open(ctx context.Context, opts terminal.Options) error {
	if opts.Command == "" {
		// Just bring Termux to front/open new session if configured in app
		return execdriver.MustRun(ctx, "am", "start", "--user", "0", "-n", "com.termux/.app.TermuxActivity").Run()
	}

	fullCmd := opts.Command
	// Resolve full path if it's just a binary name
	if !strings.HasPrefix(fullCmd, "/") {
		if path, err := execdriver.Which(ctx, fullCmd); err == nil {
			fullCmd = path
		}
	}

	if len(opts.Args) > 0 {
		// Proper escaping for the shell string
		var escapedArgs []string
		for _, arg := range opts.Args {
			escapedArgs = append(escapedArgs, fmt.Sprintf("%q", arg))
		}
		fullCmd += " " + strings.Join(escapedArgs, " ")
	}

	return execdriver.MustRun(ctx, "am", "startservice",
		"--user", "0",
		"-n", "com.termux/com.termux.app.TermuxService",
		"-a", "com.termux.service_execute",
		"-e", "com.termux.execute.command", fullCmd,
	).Run()
}
