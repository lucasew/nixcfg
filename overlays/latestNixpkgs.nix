self: super:
let
  globalConfig = import ../globalConfig.nix;
in
with super;
{
  latest = import globalConfig.latestNixpkgs {};
}
