{ pkgs, config, ... }:
{
  imports = import <dotfiles/lib/lsName.nix> ./components;

  manual.manpages.enable = false;

  home.packages =
    let
      defaultPackages =
        with pkgs; [
          fortune
          calibre
          neofetch
          file
          lazydocker
          nix-index
          scrcpy
          sqlite
          libnotify
          manix
          youtube-dl
          #browser
          google-chrome
          # compression
          xarchiver
          unzip
          p7zip
          # cloud
          rclone
          rclone-browser
          # social
          discord
          tdesktop
          # midia
          pkgs.kdeApplications.kdenlive
          pkgs.gimp
          # jetbrains
          # pkgs.jetbrains.clion

        ];
      customPackages =
        with pkgs; [
          amongUs
          usb_tixati
          # minecraft
          ets2
          mspaint
          pinball
          documentNode
          stremio
          nodePackages.vercel
        ];
      masterPackages =
        with pkgs; [
        ];
    in
    defaultPackages ++ customPackages ++ masterPackages;

  programs = {
    command-not-found.enable = true;
    jq.enable = true;
    obs-studio = {
      enable = true;
    };
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
