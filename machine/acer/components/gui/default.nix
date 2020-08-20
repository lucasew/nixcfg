let
  pkgs = import <nixpkgs> {};
in
{
  xfce = import ./xfce.nix;
  gnome = import ./gnome.nix;
  kde = import ./kde.nix;
}.${pkgs.globalConfig.selectedDesktopEnvironment}
