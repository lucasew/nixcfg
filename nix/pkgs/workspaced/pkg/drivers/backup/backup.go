package backup

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"workspaced/pkg/common"
	"workspaced/pkg/drivers/git"
	"workspaced/pkg/drivers/notification"
)

func RunFullBackup(ctx context.Context) error {
	config, err := common.LoadConfig()
	if err != nil {
		return err
	}

	logger := common.GetLogger(ctx)
	logger.Info("starting full backup")

	n := &notification.Notification{
		Title: "Backup em curso",
		Icon:  "drive-harddisk",
	}

	totalSteps := 2 // Git sync + Final report
	if common.IsRiverwood() {
		totalSteps++
	}
	if common.IsPhone() {
		totalSteps += 5 // Camera, Pictures, WA Media, WA Backups, Termux
	}

	currentStep := 0
	updateProgress := func(msg string) {
		currentStep++
		n.Message = msg
		n.Progress = float64(currentStep) / float64(totalSteps)
		n.Notify(ctx)
	}

	// Always sync git repos first
	updateProgress("Sincronizando repositórios Git...")
	git.QuickSync(ctx)

	if common.IsRiverwood() {
		updateProgress("Sincronizando CANTGIT...")
		logger.Info("host identified as riverwood, syncing CANTGIT")
		home, _ := os.UserHomeDir()
		src := filepath.Join(home, "WORKSPACE/CANTGIT/")
		dst := config.Backup.RemotePath + "/CANTGIT"
		if _, err := Rsync(ctx, src, dst, n); err != nil {
			return err
		}
	}

	if common.IsPhone() {
		logger.Info("host identified as phone, starting android backup")
		if err := runPhoneBackup(ctx, config, updateProgress, n); err != nil {
			return err
		}
	}

	// Final report
	updateProgress("Finalizando e obtendo status...")
	logger.Info("fetching remote status from rsync.net")
	status, _ := getRemoteStatus(ctx, config)
	n.Title = "Backup finalizado"
	n.Message = status
	n.Progress = 1.0
	n.Notify(ctx)

	logger.Info("full backup completed")
	return nil
}

func Rsync(ctx context.Context, src, dst string, n *notification.Notification, extraArgs ...string) (string, error) {
	config, _ := common.LoadConfig()
	remote := fmt.Sprintf("%s:%s", config.Backup.RsyncnetUser, dst)

	common.GetLogger(ctx).Info("rsync sync", "from", src, "to", remote)
	args := append([]string{"-avP", src, remote}, extraArgs...)
	cmd := common.RunCmd(ctx, "rsync", args...)

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
				n.Notify(ctx)
			}
			lastUpdate = time.Now()
		}
	}

	err = cmd.Wait()
	return lastLine, err
}

func runPhoneBackup(ctx context.Context, config *common.GlobalConfig, updateProgress func(string), n *notification.Notification) error {
	logger := common.GetLogger(ctx)
	// Sync Camera and Pictures
	logger.Info("syncing media and whatsapp")
	updateProgress("Sincronizando Câmera...")
	Rsync(ctx, "/sdcard/DCIM/Camera/", config.Backup.RemotePath+"/camera", n, "--exclude=.thumbnails")
	updateProgress("Sincronizando Fotos...")
	Rsync(ctx, "/sdcard/Pictures/", config.Backup.RemotePath+"/pictures", n, "--exclude=.thumbnails")
	updateProgress("Sincronizando Mídia WhatsApp...")
	Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Media/", config.Backup.RemotePath+"/WhatsApp", n, "--exclude=.Links", "--exclude=.Statuses")
	updateProgress("Sincronizando Backups WhatsApp...")
	Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Backups/", config.Backup.RemotePath+"/WhatsApp", n)

	// Termux config staging
	updateProgress("Sincronizando Configurações Termux...")
	logger.Info("staging termux configuration")
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".cache/backup/termux")
	os.MkdirAll(cacheDir, 0755)

	// package list
	logger.Info("generating package list")
	pkgList, _ := common.RunCmd(ctx, "dpkg-query", "-f", "${binary:Package}\n", "-W").Output()
	os.WriteFile(filepath.Join(cacheDir, "packages.txt"), pkgList, 0644)

	// sync home files
	for _, item := range []string{".bashrc", ".bash_history", ".config", ".termux", "workspace"} {
		src := filepath.Join(home, item)
		if _, err := os.Stat(src); err == nil {
			logger.Info("syncing home item", "item", item)
			common.RunCmd(ctx, "rsync", "-avP", src, cacheDir).Run()
		}
	}

	tarPath := filepath.Join(home, ".cache/backup/termux.tar")
	logger.Info("creating tarball", "path", tarPath)
	common.RunCmd(ctx, "tar", "-cvf", tarPath, "-C", filepath.Dir(cacheDir), "termux").Run()

	_, err := Rsync(ctx, tarPath, config.Backup.RemotePath, n)
	return err
}

func getRemoteStatus(ctx context.Context, config *common.GlobalConfig) (string, error) {
	user := config.Backup.RsyncnetUser

	// Get quota (raw)
	quotaOut, _ := common.RunCmd(ctx, "ssh", user, "quota").Output()

	// Get snapshots (flattened)
	snapOut, _ := common.RunCmd(ctx, "ssh", user, "ls .zfs/snapshot").Output()
	snapshots := strings.Join(strings.Fields(string(snapOut)), " ")

	return string(quotaOut) + "\n" + snapshots, nil
}

func ReplicateZFS(ctx context.Context) error {
	logger := common.GetLogger(ctx)
	// Ported from bin/misc/zfs-backup
	logger.Info("replicating ZFS vms dataset")
	if err := common.RunCmd(ctx, "sudo", "syncoid", "-r", "zroot/vms", "storage/backup/vms").Run(); err != nil {
		return err
	}
	logger.Info("replicating ZFS games dataset")
	return common.RunCmd(ctx, "sudo", "syncoid", "-r", "zroot/games", "storage/games").Run()
}
