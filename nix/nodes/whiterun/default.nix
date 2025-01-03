{
  self,
  global,
  pkgs,
  config,
  lib,
  ...
}@args:
let
  inherit (self) inputs;
  inherit (global) username;
  hostname = "whiterun";
in
{
  imports = [
    ./hardware-configuration.nix
    ../gui-common
    ./monitoring

    "${self.inputs.nixos-hardware}/common/cpu/amd/pstate.nix"
    "${self.inputs.nixos-hardware}/common/gpu/amd"
    "${self.inputs.nixos-hardware}/common/pc/ssd"
    ./nvidia.nix

    ./dlna.nix
    ./nextcloud.nix
    ./sshfs.nix
    ./zfs.nix
    ./escrivao.nix
    # ./container-inet-rdp.nix
    ./container-nat.nix
    ./rtorrent.nix
    ./transmission.nix
  ];

  services.minecraft-server.enable = true;

  services.phputils.enable = true;

  programs.nix-ld.enable = true;

  hardware.nvidia.modesetting.enable = false;

  gc-hold.enable = true;

  gc-hold.paths = with pkgs; [ blender-bin.blender_3_6 ];

  boot.kernel.sysctl = {
    "vfs.zfs.arc_sys_free" = 4 * 1024 * 1024 * 1024; # make ZFS free arc before hitting swap
  };

  services.rsyncnet-remote-backup = {
    enable = true;
    calendar = "03/6:00:01";
  };

  services.ngircd = {
    enable = true;
  };

  services.wallabag.enable = true;

  services.fusionsolar-reporter.enable = true;

  services.guix.enable = true;

  services.escrivao.enable = true;

  services.postgresql = {
    enable = true;
    enableTCPIP = true;
    userSpecificDatabases = {
      test = [ "demo" ];
    };
  };

  systemd.services.${config.services.python-microservices.services.stt-ptbr.unitName} = {
    environment = {
      LD_LIBRARY_PATH = builtins.concatStringsSep ":" [
        "/run/opengl-driver/lib"
        "${pkgs.ffmpeg.lib}/lib"
      ];
    };
  };
  services.python-microservices.services.stt-ptbr = {
    script = builtins.readFile ./python-microservice-stt-ptbr.py;
    python = pkgs.python3.withPackages (
      p: with p; [
        transformers
        torchaudio
        pytorch
      ]
    );
  };

  services.ollama.enable = true;

  networking.interfaces.enp5s0.wakeOnLan.enable = true;

  services.restic.server.enable = true;
  services.restic.server.dataDir = "/media/storage/backup/restic";

  fileSystems."/media/downloads" = {
    device = "/dev/disk/by-label/downloads";
    options = [
      "commit=60"
      "noatime"
    ];
    fsType = "ext4";
  };

  fileSystems."/media/storage" = {
    device = "/dev/disk/by-label/storage";
    fsType = "ext4";
  };

  fileSystems."/media/ssd240" = {
    device = "/dev/disk/by-label/ssd240";
    # fsType = "ext4";
  };

  fileSystems."/media/ssd1tb" = {
    device = "/dev/disk/by-label/ssd1tb";
    fsType = "ext4";
  };

  fileSystems."/var/backup" = {
    device = "/media/storage/backup/var";
    fsType = "none";
    options = [ "bind" ];
  };

  services.nginx.enable = true;

  programs.ccache.enable = true;

  services.sunshine.enable = true;

  services.xserver.windowManager.i3.enable = true;
  # programs.hyprland.enable = true;

  services.hardware.openrgb.enable = true;

  services.transmission.enable = true;

  # services.rtorrent.enable = true;
  # services.rtorrent.enableSandboxSample = true;

  services.miniflux.enable = true;

  programs.gamemode.enable = true;
  services.cf-torrent.enable = true;

  boot.plymouth.enable = true;

  services.redlib.enable = true;

  # services.boinc.enable = true;

  services.displayManager.autoLogin = {
    enable = true;
    user = "lucasew";
  };

  services.flatpak.enable = true;

  networking.hostId = "97e3b5a7";

  virtualisation.oci-containers.backend = "docker";
  virtualisation.containerd.enable = true;

  services.kubo.enable = true;

  services.telegram-sendmail.enable = true;

  services.cloud-savegame = {
    enable = true;
    calendar = "00:00:01";
    settings = {
      search = {
        paths = [
          "/media/downloads/steam/steamapps/compatdata"
          "/media/ssd240/steam/steamapps/compatdata"
        ];
      };
      flatout-2 = {
        installdir = [ "/media/downloads/steam/steamapps/common/FlatOut2" ];
      };
    };
  };

  services.nextcloud.enable = true;

  services.cockpit.enable = true;

  # services.magnetico.enable = true;

  boot = {
    supportedFilesystems = [
      "ntfs"
      "xfs"
    ];
    loader = {
      efi.canTouchEfiVariables = true;
      grub = {
        efiSupport = true;
        device = "nodev";
        # useOSProber = true; # TODO: test that VM scheme with SATA passthrough first
      };
    };
  };
  hardware.graphics.extraPackages = with pkgs; [
    rocmPackages.rocm-runtime
    rocmPackages.clr
  ];
  hardware.opengl.extraPackages32 = with pkgs; [ driversi686Linux.amdvlk ];

  programs.steam.enable = true;

  networking.hostName = hostname;

  virtualisation.libvirtd.enable = true;
  virtualisation.spiceUSBRedirection.enable = true;

  services.openssh.settings.X11Forwarding = true;

  # This value determines the NixOS release from which the default
  # settings for stateful data, like file locations and database versions
  # on your system were taken. It‘s perfectly fine and recommended to leave
  # this value at the release version of the first install of this system.
  # Before changing this value read the documentation for this option
  # (e.g. man configuration.nix or on https://nixos.org/nixos/options.html).
  system.stateVersion = "22.05"; # Did you read the comment?
}
