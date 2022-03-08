{self, global, pkgs, config, lib, ... }:
let
  hostname = "genetsec";
  inherit (self) inputs;
in {
  imports = [
    ../common/default.nix
    ./hardware-configuration.nix
    inputs.nixos-hardware.nixosModules.common-gpu-amd
    inputs.nixos-hardware.nixosModules.common-pc-hdd
  ];

  boot = {
    supportedFilesystems = [ "ntfs" ];
    loader = {
      efi.canTouchEfiVariables = true;
      grub = {
        efiSupport = true;
        device = "nodev";
        useOSProber = true; # ativado, mesmo que essa máquina n vai ver windows
      };
    };
    plymouth = {
      enable = true; #TODO: logo da utf
    };
    kernelPackages = pkgs.linuxPackages_5_10;
  };

  services.xserver = {
    enable = true;
    desktopManager = {
      xterm.enable = false;
      gnome.enable = true;
    };
    displayManager.gdm.enable = true;
    libinput.enable = true;
  };
  environment.systemPackages = with pkgs; [
    gnome3.adwaita-icon-theme
    paper-icon-theme
    p7zip zip unzip rar
    pv
    gnomeExtensions.appindicator
    gnomeExtensions.sound-output-device-chooser
  ];
  fonts.fonts = with pkgs; [
    siji
    noto-fonts
    noto-fonts-emoji
    fira-code
  ];
  programs.adb.enable = true;
  hardware.pulseaudio.enable = false;
  services = {
    auto-cpufreq.enable = true;
    gvfs.enable = true;
    printing.enable = true;

    pipewire = {
      enable = true;
      alsa.enable = true;
      alsa.support32Bit = true;
      pulse.enable = true;
      jack.enable = true;
    };
    udev.packages = with pkgs; [
      gnome3.gnome-settings-daemon
      android-udev-rules
    ];
    redshift.enable = true;
  };
  location = { # pro redshift
  latitude = -24.0;
  longitude = -54.0;
};
virtualisation = {
  docker.enable = true;
  libvirtd.enable = true;
  virtualbox.host.enable = true;
};
networking = {
  hostName = hostname;
  networkmanager.enable = true;
  defaultGateway = "192.168.100.2";
  nameservers = [ "8.8.8.8" "8.8.4.4" ];
  interfaces.enp3s0 = {
    useDHCP = false;
    subnetMask = "255.255.255.0";
    ipv4.addresses = [
      { address = "192.168.100.56"; prefixLength = 24; }
    ];
  };
};
users.users = {
  genetsec = {
    extraGroups = [
      "adbusers"
      "vboxusers"
    ];
    description = "GENETSEC";
    isNormalUser = true;
    initialPassword = "genetsec";
  };
};
  # Esse default timeout é pra quando o pc ta iniciando ele não ficar esperando muito pra serviços iniciarem, já teve bug de serviço travar boot por causa de não ter internet
  systemd.extraConfig = ''
  DefaultTimeoutStartSec=10s
  '';
  systemd.services.NetworkManager-wait-online.enable = false;
  hardware = {
    bluetooth.enable = true;
    opengl = {
      enable = true;
      driSupport32Bit = true;
      extraPackages32 = with pkgs.pkgsi686Linux; [
        vaapiIntel
      ];
    };
  };

  fileSystems = {
    "/" = {
      device = "/dev/disk/by-label/nixos";
      fsType = "ext4";
    };
    "/boot" = {
      device = "/dev/disk/by-label/ESP";
      fsType = "vfat";
    };
  };
  swapDevices = [
    { device = "/dev/disk/by-label/swap"; }
  ];

  system.stateVersion = "20.03";
}

