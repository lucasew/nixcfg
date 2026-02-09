package backup

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"workspaced/pkg/config"
	"workspaced/pkg/env"
	"workspaced/pkg/exec"
	"workspaced/pkg/git"
	"workspaced/pkg/logging"
	"workspaced/pkg/notification"
	"workspaced/pkg/sudo"
	"workspaced/pkg/types"
)

func RunFullBackup(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	logger := logging.GetLogger(ctx)
	logger.Info("starting full backup")

	n := &notification.Notification{
		Title: "Backup em curso",
		Icon:  "drive-harddisk",
	}

	totalSteps := 2 // Git sync + Final report
	if env.IsRiverwood() {
		totalSteps++
	}
	if env.IsPhone() {
		totalSteps += 5 // Camera, Pictures, WA Media, WA Backups, Termux
	}

	currentStep := 0
	updateProgress := func(msg string) {
		currentStep++
		n.Message = msg
		n.Progress = float64(currentStep) / float64(totalSteps)
		_ = n.Notify(ctx)
	}

	// Always sync git repos first
	updateProgress("Sincronizando repositórios Git...")
	_ = git.QuickSync(ctx)

	if env.IsRiverwood() {
		updateProgress("Sincronizando CANTGIT...")
		logger.Info("host identified as riverwood, syncing CANTGIT")
		home, _ := os.UserHomeDir()
		src := filepath.Join(home, "WORKSPACE/CANTGIT/")
		dst := cfg.Backup.RemotePath + "/CANTGIT"
		if _, err := Rsync(ctx, src, dst, n); err != nil {
			return err
		}
	}

	if env.IsPhone() {
		logger.Info("host identified as phone, starting android backup")
		if err := runPhoneBackup(ctx, cfg, updateProgress, n); err != nil {
			return err
		}
	}

	// Final report
	updateProgress("Finalizando e obtendo status...")
	logger.Info("fetching remote status from rsync.net")
	status, _ := getRemoteStatus(ctx, cfg)
	n.Title = "Backup finalizado"
	n.Message = status
	n.Progress = 1.0
	_ = n.Notify(ctx)

	logger.Info("full backup completed")
	return nil
}

func Rsync(ctx context.Context, src, dst string, n *notification.Notification, extraArgs ...string) (string, error) {
	cfg, _ := config.LoadConfig()
	remote := fmt.Sprintf("%s:%s", cfg.Backup.RsyncnetUser, dst)

	logging.GetLogger(ctx).Info("rsync sync", "from", src, "to", remote)
	args := append([]string{"-avP", src, remote}, extraArgs...)
	cmd := exec.RunCmd(ctx, "rsync", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return "", err
	}

	lastLine := ""
	scanner := bufio.NewScanner(stdout)
	lastUpdate := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lastLine = line
		}
		if time.Since(lastUpdate) > time.Second {
			if n != nil {
				n.Message = line
				_ = n.Notify(ctx)
			}
			lastUpdate = time.Now()
		}
	}

	err = cmd.Wait()
	return lastLine, err
}

func runPhoneBackup(ctx context.Context, cfg *config.GlobalConfig, updateProgress func(string), n *notification.Notification) error {
	logger := logging.GetLogger(ctx)
	// Sync Camera and Pictures
	logger.Info("syncing media and whatsapp")
	updateProgress("Sincronizando Câmera...")
	_, _ = Rsync(ctx, "/sdcard/DCIM/Camera/", cfg.Backup.RemotePath+"/camera", n, "--exclude=.thumbnails")
	updateProgress("Sincronizando Fotos...")
	_, _ = Rsync(ctx, "/sdcard/Pictures/", cfg.Backup.RemotePath+"/pictures", n, "--exclude=.thumbnails")
	updateProgress("Sincronizando Mídia WhatsApp...")
	_, _ = Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Media/", cfg.Backup.RemotePath+"/WhatsApp", n, "--exclude=.Links", "--exclude=.Statuses")
	updateProgress("Sincronizando Backups WhatsApp...")
	_, _ = Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Backups/", cfg.Backup.RemotePath+"/WhatsApp", n)

	// Termux config staging
	updateProgress("Sincronizando Configurações Termux...")
	logger.Info("staging termux configuration")
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".cache/backup/termux")
	_ = os.MkdirAll(cacheDir, 0755)

	// package list
	logger.Info("generating package list")
	pkgList, _ := exec.RunCmd(ctx, "dpkg-query", "-f", "${binary:Package}\n", "-W").Output()
	_ = os.WriteFile(filepath.Join(cacheDir, "packages.txt"), pkgList, 0644)

	// sync home files
	for _, item := range []string{".bashrc", ".bash_history", ".config", ".termux", "workspace"} {
		src := filepath.Join(home, item)
		if _, err := os.Stat(src); err == nil {
			logger.Info("syncing home item", "item", item)
			_ = exec.RunCmd(ctx, "rsync", "-avP", src, cacheDir).Run()
		}
	}

	tarPath := filepath.Join(home, ".cache/backup/termux.tar")
	logger.Info("creating tarball", "path", tarPath)
	_ = exec.RunCmd(ctx, "tar", "-cvf", tarPath, "-C", filepath.Dir(cacheDir), "termux").Run()

	_, err := Rsync(ctx, tarPath, cfg.Backup.RemotePath, n)
	return err
}

func getRemoteStatus(ctx context.Context, cfg *config.GlobalConfig) (string, error) {
	user := cfg.Backup.RsyncnetUser

	// Get quota (raw)
	quotaOut, _ := exec.RunCmd(ctx, "ssh", user, "quota").Output()

	// Filter out lines with asterisks from quota output
	var quotaLines []string
	for _, line := range strings.Split(string(quotaOut), "\n") {
		if !strings.Contains(line, "*") && line != "" {
			quotaLines = append(quotaLines, line)
		}
	}
	filteredQuota := strings.Join(quotaLines, "\n")

	// Get snapshots (flattened)
	snapOut, _ := exec.RunCmd(ctx, "ssh", user, "ls .zfs/snapshot").Output()
	snapshots := strings.Join(strings.Fields(string(snapOut)), " ")

	return filteredQuota + "\n" + snapshots, nil
}

func ReplicateZFS(ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	// Ported from bin/misc/zfs-backup
	logger.Info("replicating ZFS vms dataset")

	if os.Getuid() != 0 {
		_ = sudo.Enqueue(ctx, &types.SudoCommand{
			Slug:    "zfs-backup-vms",
			Command: "syncoid",
			Args:    []string{"-r", "zroot/vms", "storage/backup/vms"},
		})
		logger.Info("replicating ZFS games dataset")
		_ = sudo.Enqueue(ctx, &types.SudoCommand{
			Slug:    "zfs-backup-games",
			Command: "syncoid",
			Args:    []string{"-r", "zroot/games", "storage/games"},
		})
		return nil
	}

	if err := exec.RunCmd(ctx, "syncoid", "-r", "zroot/vms", "storage/backup/vms").Run(); err != nil {
		return err
	}
	logger.Info("replicating ZFS games dataset")
	return exec.RunCmd(ctx, "syncoid", "-r", "zroot/games", "storage/games").Run()
}
