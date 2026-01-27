package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
		Title:   "Backup Iniciado",
		Message: "Sincronizando dados...",
		Icon:    "drive-harddisk",
	}
	n.Notify(ctx)

	// Always sync git repos first
	git.QuickSync(ctx)

	if common.IsRiverwood() {
		logger.Info("host identified as riverwood, syncing CANTGIT")
		home, _ := os.UserHomeDir()
		src := filepath.Join(home, "WORKSPACE/CANTGIT/")
		dst := config.Backup.RemotePath + "/CANTGIT"
		if err := Rsync(ctx, src, dst); err != nil {
			return err
		}
	}

	if common.IsPhone() {
		logger.Info("host identified as phone, starting android backup")
		if err := runPhoneBackup(ctx, config); err != nil {
			return err
		}
	}

	// Final report
	logger.Info("fetching remote status from rsync.net")
	status, _ := getRemoteStatus(ctx, config)
	n.Message = "Backup finalizado.\n" + status
	n.Notify(ctx)

	logger.Info("full backup completed")
	return nil
}

func Rsync(ctx context.Context, src, dst string, extraArgs ...string) error {
	config, _ := common.LoadConfig()
	remote := fmt.Sprintf("%s:%s", config.Backup.RsyncnetUser, dst)

	common.GetLogger(ctx).Info("rsync sync", "from", src, "to", remote)
	args := append([]string{"-avP", src, remote}, extraArgs...)
	return common.RunCmd(ctx, "rsync", args...).Run()
}

func runPhoneBackup(ctx context.Context, config *common.GlobalConfig) error {
	logger := common.GetLogger(ctx)
	// Sync Camera and Pictures
	logger.Info("syncing media and whatsapp")
	Rsync(ctx, "/sdcard/DCIM/Camera/", config.Backup.RemotePath+"/camera", "--exclude=.thumbnails")
	Rsync(ctx, "/sdcard/Pictures/", config.Backup.RemotePath+"/pictures", "--exclude=.thumbnails")
	Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Media/", config.Backup.RemotePath+"/WhatsApp", "--exclude=.Links", "--exclude=.Statuses")
	Rsync(ctx, "/sdcard/Android/media/com.whatsapp/WhatsApp/Backups/", config.Backup.RemotePath+"/WhatsApp")

	// Termux config staging
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

	return Rsync(ctx, tarPath, config.Backup.RemotePath)
}

func getRemoteStatus(ctx context.Context, config *common.GlobalConfig) (string, error) {
	user := config.Backup.RsyncnetUser

	// Get snapshots
	snapOut, _ := common.RunCmd(ctx, "ssh", user, "ls .zfs/snapshot").Output()
	// Get quota
	quotaOut, _ := common.RunCmd(ctx, "ssh", user, "quota").Output()

	return string(snapOut) + "\n" + string(quotaOut), nil
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
