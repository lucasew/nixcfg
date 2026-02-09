package sudo

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
	"workspaced/pkg/types"
)

func getQueueDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cache/workspaced/sudo_queue")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func Enqueue(ctx context.Context, cmd *types.SudoCommand) error {
	if cmd.Slug == "" {
		b := make([]byte, 3)
		if _, err := rand.Read(b); err != nil {
			return err
		}
		cmd.Slug = fmt.Sprintf("%x", b)
	}

	if cmd.Cwd == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.Cwd = cwd
	}

	if len(cmd.Env) == 0 {
		cmd.Env = os.Environ()
	}

	if cmd.Timestamp == 0 {
		cmd.Timestamp = time.Now().Unix()
	}

	dir, err := getQueueDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, cmd.Slug+".json")
	data, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	n := &notification.Notification{
		Title:   "Sudo Required",
		Message: fmt.Sprintf("Command '%s' (slug: %s) pending approval.", cmd.Command, cmd.Slug),
		Icon:    "dialog-password",
	}
	if err := n.Notify(ctx); err != nil {
		logging.ReportError(ctx, err)
	}

	return nil
}

func List(ctx context.Context) ([]*types.SudoCommand, error) {
	dir, err := getQueueDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var cmds []*types.SudoCommand
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				logging.ReportError(ctx, err)
				continue
			}
			var cmd types.SudoCommand
			if err := json.Unmarshal(data, &cmd); err != nil {
				logging.ReportError(ctx, err)
				continue
			}
			cmds = append(cmds, &cmd)
		}
	}
	return cmds, nil
}

func Get(slug string) (*types.SudoCommand, error) {
	dir, err := getQueueDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, slug+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cmd types.SudoCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func Remove(slug string) error {
	dir, err := getQueueDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, slug+".json")
	return os.Remove(path)
}
