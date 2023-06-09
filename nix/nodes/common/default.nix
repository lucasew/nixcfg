{self, pkgs, lib, unpackedInputs, ...}:
let
  inherit (self) inputs;
  inherit (lib) mkDefault;
in
{
  imports = [
    ./nginx.nix
    ./nginx-root-domain.nix
    ./nix-index-database.nix
    ./netusage.nix
    ./magnetico.nix
    ./lvm.nix
    ./transmission
    ./invidious.nix
    ./cockpit-extra.nix
    ./libreddit.nix
    ../bootstrap/default.nix
    ../../modules/cachix/system.nix
    ../../modules/hold-gc/system.nix
    ./ansible-python.nix
    ./cloud-savegame.nix
    ./hosts.nix
    ./p2k.nix
    ./tuning.nix
    ./cf-torrent.nix
    ./tt-rss.nix
    ./jellyfin.nix
    ./user.nix
    ./tmux
    ./bash
    ./kvm.nix
    ./sops.nix
    ./unstore.nix
    ./telegram_sendmail.nix
    "${unpackedInputs.simple-dashboard}/nixos-module.nix"
  ];

  services.lvm.enable = mkDefault false;

  programs.fuse.userAllowOther = true;

  services.cloud-savegame = {
    enableVerbose = true;
    enableGit = true;
    settings = {
      search = {
        paths = [ "~" ];
        extra_homes = [ "/run/media/lucasew/Dados/DADOS/Lucas" ];
      };

      flatout-2 = {
        installdir= [ "~/.local/share/Steam/steamapps/common/FlatOut2" "/run/media/lucasew/Dados/DADOS/Jogos/FlatOut 2"];
      };

      farming-simulator-2013 = {
        ignore_mods = true;
      };
    };
  };



  services.unstore = {
    # enable = true;
    paths = [
      "flake.nix"
    ];
  };

  boot.loader.grub.memtest86.enable = true;

  virtualisation.docker = {
    enable = true;
    # dockerSocket.enable = true;
    # dockerCompat = true;
    enableNvidia = true;
  };

  environment = {
    systemPackages = with pkgs; [
      rlwrap
      wget
      curl
      unrar
      direnv
      pciutils
      usbutils
      htop
      lm_sensors
      neofetch
    ];
  };
  cachix.enable = true;

  services.smartd = {
    enable = true;
    autodetect = true;
    notifications.test = true;
  };

  services.nginx.appendHttpConfig = ''
    error_log stderr;
    access_log syslog:server=unix:/dev/log combined;
  '';
}
