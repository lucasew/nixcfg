{ lib, ... }:
{
  imports = [
    ./cloud-savegame.nix
    ./cockpit-extra.nix
    ./cf-torrent.nix
    ./fusionsolar
    ./magnetico.nix
    ./redlib.nix
    ./transmission.nix
    ./telegram_sendmail.nix
    ./postgresql.nix
    ./nixgram.nix
    ./netusage
    ./restic-server.nix
    ./python-microservices
    ./rtorrent.nix
    ./rsyncnet
    ./wallabag.nix
    ./ts-proxy.nix
    ./ngircd.nix
    ./minecraft.nix
    ./phpelo.nix
  ];

  virtualisation.oci-containers.backend = lib.mkDefault "podman";
}
