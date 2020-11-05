self: super:
let
  npmPackages = import ./package_data/default.nix;
in
{
  nodePackages = super.nodePackages // npmPackages {};
}
