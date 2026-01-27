package main

import (
	"bytes"
	"context"

	"github.com/spf13/cobra"
)

func runCobra(root *cobra.Command, args []string, env []string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	ctx := context.WithValue(context.Background(), "env", env)

	err := root.ExecuteContext(ctx)
	return buf.String(), err
}
