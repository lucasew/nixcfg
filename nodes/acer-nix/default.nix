# Edit this configuration file to define what should be installed on
# your system.  Help is available in the configuration.nix(5) man page
# and in the NixOS manual (accessible by running ‘nixos-help’).

{pkgs, config, nix-ld, ... }:
let
  username = "lucasew";
  hostname = "acer-nix";
in
{
  imports =
    [
      # Include the results of the hardware scan.
      ./hardware-configuration.nix
    ]
    ++ [
      ./modules/virt-manager/system.nix
      ../../modules/cachix/system.nix
      ../../modules/gui/system.nix
      ../../modules/polybar/system.nix
    ]
  ;
  # gui = {
  #   enable = true;
  #   selected = "xfce_i3";
  # };

  nixpkgs.config.allowUnfree = true;
  nix = {
    package = pkgs.nixFlakes;
    autoOptimiseStore = true;
    gc = {
      options = "--delete-older-than 15d";
    };
    extraOptions = ''
      min-free = ${toString (1  *1024*1024*1024)}
      max-free = ${toString (10 *1024*1024*1024)}
      experimental-features = nix-command flakes
    '';
  };

  # Use the systemd-boot EFI boot loader.
  boot.supportedFilesystems = [ "ntfs" ];
  boot.loader = {
    efi = {
      canTouchEfiVariables = true;
    };
    grub = {
      efiSupport = true;
      #efiInstallAsRemovable = true; # in case canTouchEfiVariables doesn't work for your system
      device = "nodev";
      useOSProber = true;
    };
  };

  # limpar tmp no boot
  boot.cleanTmpDir = true;

  networking.hostName = hostname; # Define your hostname.
  networking.networkmanager.enable = true;

  # The global useDHCP flag is deprecated, therefore explicitly set to false here.
  # Per-interface useDHCP will be mandatory in the future, so this generated config
  # replicates the default behaviour.
  networking.useDHCP = false;
  networking.interfaces.enp2s0f1.useDHCP = true;
  networking.interfaces.wlp3s0.useDHCP = true;

  # Configure network proxy if necessary
  # networking.proxy.default = "http://user:password@proxy:port/";
  # networking.proxy.noProxy = "127.0.0.1,localhost,internal.domain";
  # Select internationalisation properties.
  i18n.defaultLocale = "pt_BR.UTF-8";
  # console = {
  #   font = "Lat2-Terminus16";
  #   keyMap = "us";
  # };

  # Set your time zone.
  time.timeZone = "America/Sao_Paulo";

  # List packages installed in system profile. To search, run:
  # $ nix search wget
  environment.systemPackages = with pkgs; [
    wget
    gparted
    paper-icon-theme
    kde-gtk-config # Custom
    dasel # manipulação de json, toml, yaml, xml, csv e tal
    rclone rclone-browser restic # cloud storage
    p7zip unzip xarchiver # archiving
    (pkgs.callPackage ../../modules/neovim/package.nix {})
    # Extra
    gitAndTools.gitui
  ];

  # melhor editor ever
  environment.variables.EDITOR = "nvim";

  programs.dconf.enable = true;
  services.dbus.packages = with pkgs; [ gnome3.dconf ];

  # Some programs need SUID wrappers, can be configured further or are
  # started in user sessions.
  # programs.mtr.enable = true;
  # programs.gnupg.agent = {
  #   enable = true;
  #   enableSSHSupport = true;
  #   pinentryFlavor = "gnome3";
  # };
  hardware.opengl.enable = true;
  hardware.opengl.driSupport32Bit = true;

  # Enable the OpenSSH daemon.
  services.openssh.enable = true;

  # Open ports in the firewall.
  # networking.firewall.allowedTCPPorts = [ ... ];
  # networking.firewall.allowedUDPPorts = [ ... ];
  # Or disable the firewall altogether.
  networking.firewall.enable = false;

  # Enable CUPS to print documents.
  # services.printing.enable = true;

  # Enable sound.
  sound.enable = true;
  hardware.pulseaudio.enable = true;

  # services.xserver.layout = "us";
  # services.xserver.xkbOptions = "eurosign:e";

  # Enable touchpad support.
  services.xserver.libinput.enable = true;

  # Tailscale
  services.tailscale.enable = true;

  # Themes
  programs.qt5ct.enable = true;

  # Users
  users.users = {
    ${username} = {
      isNormalUser = true;
      extraGroups = [
        "wheel" # sudo
        "docker" # docker
        "adbusers" # adb
      ]; 
      description = "Lucas Eduardo";
    };
  };
  # Home manager
  home-manager = {
    users = {
      "${username}" = import ./home.nix;
    };
    useUserPackages = true;
#    useGlobalPkgs = true;
  };
  # ADB
  programs.adb.enable = true;
  services.udev.packages = with pkgs; [
    gnome3.gnome-settings-daemon
    android-udev-rules
  ];
  # docker
  virtualisation.docker.enable = true;

  # keybase
  services = {
    keybase.enable = true;
    kbfs.enable = true;
  };

  # cachix
  cachix.enable = true;

  # singularity
  programs.singularity.enable = true;

  # não deixar explodir
  nix.maxJobs = 3;
  # kernel
  boot.kernelPackages = pkgs.linuxPackages_5_10;

  # This value determines the NixOS release from which the default
  # settings for stateful data, like file locations and database versions
  # on your system were taken. It‘s perfectly fine and recommended to leave
  # this value at the release version of the first install of this system.
  # Before changing this value read the documentation for this option
  # (e.g. man configuration.nix or on https://nixos.org/nixos/options.html).
  system.stateVersion = "20.03"; # Did you read the comment?
}
