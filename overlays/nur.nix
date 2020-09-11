self: super:
let
  pkgs = super.pkgs;
  globalConfig = import <dotfiles/globalConfig.nix>;
  nurRepo = globalConfig.nur;
in
{
  nur = import nurRepo {
    inherit pkgs;
  };
}
