{lib, ...}: {
  imports = [
    ./cloud-savegame.nix
    ./cockpit-extra.nix
    ./transmission.nix
    ./telegram_sendmail.nix
    ./nomad.nix
    ./netusage
    ./restic-server.nix
    ./python-microservices
    ./rtorrent.nix
    ./ts-proxy.nix
    ./ngircd.nix
    ./minecraft.nix
    ./phpelo.nix
  ];

  virtualisation.oci-containers.backend = lib.mkDefault "podman";
}
