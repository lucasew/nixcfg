let
    overlays = import <dotfiles/overlays/utils/importAllIn.nix> ./overlays;
in import <nixpkgs> {overlays = overlays;}