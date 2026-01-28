package nix

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/nix"
	"workspaced/pkg/drivers/notification"

	"github.com/spf13/cobra"
)

func init() {
	Registry.Register(func(parent *cobra.Command) {
		cmd := &cobra.Command{
			Use:   "deploy [nodes...]",
			Short: "Deploy NixOS and Home Manager configurations to remote nodes",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := cmd.Context()
				nodes := args
				if len(nodes) == 0 {
					nodes = []string{"riverwood", "whiterun"}
				}

				flake, _ := cmd.Flags().GetString("flake")
				if flake == "" {
					root, err := common.GetDotfilesRoot()
					if err != nil {
						return err
					}
					flake = root
				}

				action, _ := cmd.Flags().GetString("action")

				for _, node := range nodes {
					slog.Info("Deploying to node", "node", node)
					if err := deployNode(ctx, flake, node, action); err != nil {
						slog.Error("Failed to deploy to node", "node", node, "error", err)
						return err
					}
				}

				n := &notification.Notification{
					Title:   "NixOS Deploy",
					Message: fmt.Sprintf("Deploy concluído para: %s", strings.Join(nodes, ", ")),
					Icon:    "nix-snowflake",
				}
				n.Notify(ctx)

				return nil
			},
		}
		cmd.Flags().StringP("flake", "f", "", "Flake reference to use")
		cmd.Flags().StringP("action", "a", "", "Action to perform (switch, boot, test). If empty, auto-detects.")
		parent.AddCommand(cmd)
	})
}

func deployNode(ctx context.Context, flake, node, action string) error {
	// 1. Build outputs
	slog.Info("Building configuration for node", "node", node)
	toplevelPath := fmt.Sprintf("nixosConfigurations.%s.config.system.build.toplevel", node)
	toplevel, err := nix.GetFlakeOutput(ctx, flake, toplevelPath)
	if err != nil {
		return fmt.Errorf("failed to build toplevel for %s: %w", node, err)
	}

	homePath := "homeConfigurations.main.activationPackage"
	home, err := nix.GetFlakeOutput(ctx, flake, homePath)
	if err != nil {
		return fmt.Errorf("failed to build home-manager for %s: %w", node, err)
	}

	// 2. Copy closures
	slog.Info("Copying closures to node", "node", node)
	if err := nix.CopyClosure(ctx, node, toplevel, nix.To); err != nil {
		return fmt.Errorf("failed to copy toplevel to %s: %w", node, err)
	}
	if err := nix.CopyClosure(ctx, node, home, nix.To); err != nil {
		return fmt.Errorf("failed to copy home-manager to %s: %w", node, err)
	}

	// 3. Auto-detect action if not specified
	if action == "" {
		action = "boot"
		// Check if same nixpkgs is used
		localUsedOut, err := common.RunCmd(ctx, "realpath", fmt.Sprintf("%s/etc/.nixpkgs-used", toplevel)).Output()
		if err == nil {
			localUsed := strings.TrimSpace(string(localUsedOut))
			remoteUsedOut, err := common.RunCmd(ctx, "ssh", node, "realpath /etc/.nixpkgs-used").Output()
			if err == nil {
				remoteUsed := strings.TrimSpace(string(remoteUsedOut))
				if localUsed == remoteUsed {
					action = "switch"
				}
			}
		}
	}

	// 4. Activate Home Manager
	slog.Info("Activating Home Manager on node", "node", node)
	cmdHM := common.RunCmd(ctx, "ssh", node, fmt.Sprintf("%s/bin/home-manager-generation", home))
	common.InheritContextWriters(ctx, cmdHM)
	if err := cmdHM.Run(); err != nil {
		return fmt.Errorf("failed to activate home-manager on %s: %w", node, err)
	}

	// 5. Switch System Configuration
	slog.Info("Switching system configuration on node", "node", node, "action", action)
	// Check if already running
	currentSystemOut, err := common.RunCmd(ctx, "ssh", node, "realpath /run/current-system").Output()
	if err == nil {
		currentSystem := strings.TrimSpace(string(currentSystemOut))
		if currentSystem == toplevel {
			slog.Info("Node already running the same configuration", "node", node)
			return nil
		}
	}

	switchCmdArgs := []string{"ssh", "-t", node, "sudo", fmt.Sprintf("%s/bin/switch-to-configuration", toplevel), action}
	cmdSwitch := common.RunCmd(ctx, switchCmdArgs[0], switchCmdArgs[1:]...)
	common.InheritContextWriters(ctx, cmdSwitch)
	if err := cmdSwitch.Run(); err != nil {
		return fmt.Errorf("failed to switch configuration on %s: %w", node, err)
	}

	return nil
}
