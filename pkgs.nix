let
  globalConfig = import ./globalConfig.nix;
  overlays = import <dotfiles/overlays/utils/importAllIn.nix> globalConfig.overlaysPath;
  nixpkgs = import globalConfig.nixpkgs;
in nixpkgs {overlays = overlays;}
