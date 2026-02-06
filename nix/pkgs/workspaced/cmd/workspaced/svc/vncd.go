package svc

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"syscall"
	"workspaced/pkg/exec"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		parent.AddCommand(&cobra.Command{
			Use:   "vncd",
			Short: "Start a VNC server (Wayland or X11)",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				waylandDisplay := os.Getenv("WAYLAND_DISPLAY")

				if waylandDisplay != "" {
					return runWaylandVNC(ctx)
				}
				return runXorgVNC(ctx)
			},
		})
	})
}

func runWaylandVNC(ctx context.Context) error {
	slog.Info("Starting wayvnc")
	host := os.Getenv("WAYVNC_HOST")
	if host == "" {
		tsIP, err := getTailscaleIP(ctx)
		if err == nil && tsIP != "" {
			host = tsIP
		} else {
			host = "127.0.0.1"
		}
	}

	bin, err := exec.Which(ctx, "wayvnc")
	if err != nil {
		return err
	}

	slog.Info("executing wayvnc", "host", host)
	return syscall.Exec(bin, []string{"wayvnc", host}, os.Environ())
}

func runXorgVNC(ctx context.Context) error {
	slog.Info("Starting x0vncserver")
	bin, err := exec.Which(ctx, "x0vncserver")
	if err != nil {
		return err
	}

	args := []string{
		"x0vncserver",
		"-display=:0",
		"-SecurityTypes", "None",
		"-ImprovedHextile=1",
		"-RawKeyboard=1",
	}

	slog.Info("executing x0vncserver")
	return syscall.Exec(bin, args, os.Environ())
}

func getTailscaleIP(ctx context.Context) (string, error) {
	out, err := exec.RunCmd(ctx, "tailscale", "ip", "-4").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
