{
  global,
  pkgs,
  lib,
  ...
}:

{

  imports = [
    ../base/default.nix
    ./ghostty.nix
    ./dconf.nix
    ./theme
  ];

  stylix.enable = true;

  programs.man.enable = true;

  home = {
    homeDirectory = /home/lucasew;
    inherit (global) username;
  };

  home.packages = with pkgs; [
    mission-center
    nix-output-monitor
    pkg

    # media
    nbr.wine-apps._7zip
    xxd

    # LSPs
    nil
    python3Packages.python-lsp-server

    (pkgs.makeDesktopItem {
      name = "nixcfg-quicksync";
      desktopName = "nixcfg: Sincronização Rápida";
      icon = "sync-synchronizing";
      exec = "sdw quicksync";
    })
    (pkgs.makeDesktopItem {
      name = "nixcfg-backup";
      desktopName = "nixcfg: Backup";
      icon = "sync-synchronizing";
      exec = "sdw backup";
    })
  ];

  gtk = {
    enable = true;
  };
  qt = {
    enable = true;
    platformTheme.name = "gtk";
  };

  programs.terminator = {
    # enable = true;
    config = {
      global_config.borderless = true;
    };
  };
}
