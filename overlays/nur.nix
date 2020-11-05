self: super:
with super;
let
  globalConfig = import <dotfiles/globalConfig.nix>;
  nurRepo = globalConfig.nur;
in
{
  nur = import nurRepo {
    inherit pkgs;
  };
}
