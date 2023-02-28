{ ... }:

# syncing to backup zpool
#   syncoid storage/vmiso archive/downloads/vmiso
#   syncoid storage/backup-hdexterno archive/backup-hdexterno
#   syncoid zroot/vms archive/vms -r

{
  boot.supportedFilesystems = [ "zfs" ];
  services.zfs.autoScrub = {
    enable = true;
    pools = [ "storage" "zroot" ];
  };
  boot.zfs.requestEncryptionCredentials = [ "zroot" ];
  boot.zfs.extraPools = [ "storage" ];
  virtualisation.docker.storageDriver = "zfs";
}
