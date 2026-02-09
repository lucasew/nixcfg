{ pkgs, ... }:
{
  services.xserver = {
    desktopManager.xterm.enable = false;
  };

  environment.systemPackages = with pkgs; [
    lxappearance
    glib # gsettings
    gsettings-desktop-schemas # schemas necessários para GTK
  ];

  fonts.packages = with pkgs; [
    siji
    noto-fonts
    noto-fonts-color-emoji
    fira-code
  ];

  services.xserver = {
    xkb = {
      layout = "br,us";
      options = "grp:win_space_toggle,terminate:ctrl_alt_bksp";
      variant = ",";
    };
  };

  # Enable touchpad support.
  services.libinput.enable = true;

  location = {
    latitude = -24.0;
    longitude = -54.0;
  };

  # dconf para persistir configurações GTK (lxappearance, etc)
  programs.dconf.enable = true;

  # Configurar gsettings schemas
  services.dbus.packages = with pkgs; [
    dconf
    gsettings-desktop-schemas
  ];
}
