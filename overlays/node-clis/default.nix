self: super:
let
  npmPackages = import ./assets/default.nix;
in
{
  nodePackages = super.nodePackages // npmPackages {};
}
