{ config, ... }:
let
  pkgs = import <dotfiles/pkgs.nix>;
  bin = pkgs.writeShellScriptBin "nix-search" ''
    ${pkgs.dotwrap}/bin/dotwrap nix search -f '<dotfiles/pkgs.nix>' $*
  '';
in
{
  home.packages = [
    bin
  ];
}
