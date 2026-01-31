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

  home = {
    homeDirectory = /home/lucasew;
    inherit (global) username;
  };

  home.packages = with pkgs; [
    pkg

    # media
    nbr.wine-apps._7zip

    # LSPs
    nil
    python3Packages.python-lsp-server
  ];

  gtk = {
    enable = true;
  };
  qt = {
    enable = true;
    platformTheme.name = "gtk";
  };
}
