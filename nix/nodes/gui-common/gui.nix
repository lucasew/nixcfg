{ pkgs, ... }:
{
  services.xserver = {
    desktopManager.xterm.enable = false;
  };

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
}
