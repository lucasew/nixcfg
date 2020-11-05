{pkgs, ...}:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
  nurRepo = globalConfig.nur;
in import nurRepo {
  inherit pkgs;
}
