{
  pkgs,
  lib,
  global,
  ...
}:
let
  inherit (global) username;
in
{
  imports = [
    ../common
    ./gui-variants
    ./audio.nix
    ./gui.nix
    ./networking.nix
    ./steam.nix
    ./git.nix
    ./gammastep.nix
    ./adb.nix
    ./vbox.nix
    ./tuning.nix
    ./ipfs.nix
    ./gamemode.nix
    ./sunshine.nix
    ./wallpaper.nix
    ./extra-fonts.nix
    ./polkit.nix
    ./gui-variants
  ];

  systemd.extraConfig = ''
    DefaultTimeoutStartSec=10s
  '';

  environment.systemPackages = with pkgs; [
    keepassxc
    parallel
    home-manager
    paper-icon-theme
    p7zip
    unzip # archiving
    pv
    # Extra
    distrobox # plan b
    xorg.xkill
  ];

  programs.dconf.enable = true;
  services.dbus.packages = with pkgs; [ dconf ];
  services.gvfs.enable = true;
  services.tumbler.enable = true;

  programs.ssh = {
    startAgent = true;
    extraConfig = ''
      ConnectTimeout=5
    '';
  };
  services.shellhub-agent = {
    enable = true;
    tenantId = "c574bf33-a21a-49ef-a7a5-1d8fbd823e4e";
  };
  programs.gnupg.agent = {
    enable = true;
    # enableSSHSupport = true;
    pinentryPackage = pkgs.pinentry-gnome3;
  };

  # Users
  users.users = {
    ${username} = {
      description = "Lucas Eduardo";
    };
  };

  hardware.graphics.enable = true;

  # Enable CUPS to print documents.
  services.printing.enable = lib.mkDefault true;

  qt.platformTheme.name = lib.mkDefault "qt5ct";

  # https://github.com/NixOS/nixpkgs/pull/297434#issuecomment-2348783988
  systemd.services.display-manager.environment.XDG_CURRENT_DESKTOP = "X-NIXOS-SYSTEMD-AWARE";
}
