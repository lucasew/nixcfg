let
  globalConfig = import <dotfiles/globalConfig.nix>;
in import (../gui + "/${globalConfig.selectedDesktopEnvironment}.nix")
