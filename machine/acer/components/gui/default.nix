{
  xfce = import ./xfce.nix;
  gnome = import ./gnome.nix;
  kde = import ./kde.nix;
}.${import ../../../../config.nix "selectedDesktopEnvironment"}
