{ config, pkgs, ... }:
{
  imports = (import ../overlays/utils/lsName.nix) ./components;

  home.packages = with pkgs; [
    fortune
    calibre
    tdesktop
    neofetch
    file
    arduino
    heroku
    lazydocker
    youtube-dl
    nix-index
    scrcpy
    sqlite
    libnotify
  ];

  programs.command-not-found.enable = true;
  programs.jq.enable = true;
  programs.obs-studio = {
    enable = true;
    plugins = [];
  };
  # This value determines the Home Manager release that your
  # configuration is compatible with. This helps avoid breakage
  # when a new Home Manager release introduces backwards
  # incompatible changes.
  #
  # You can update Home Manager without changing this value. See
  # the Home Manager release notes for a list of state version
  # changes in each release.
  gtk = {
    enable = true;
    theme.name = "Adwaita-dark";
  };
  qt = {
    enable = true;
    platformTheme = "gtk";
  };

  home.stateVersion = "20.03";

}
