let
  pkgs = import <dotfiles/pkgs.nix>;
in
{
  xfce = import ./xfce.nix;
  gnome = import ./gnome.nix;
  kde = import ./kde.nix;
}.${pkgs.globalConfig.selectedDesktopEnvironment}
