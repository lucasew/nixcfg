let
  pkgs = import <dotfiles/pkgs.nix>;
in import (../gui + "/${pkgs.globalConfig.selectedDesktopEnvironment}.nix")
