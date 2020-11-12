self: super:
let
  npmPackages = import ./package_data/default.nix {pkgs = super.pkgs;};
  a22120 = npmPackages."22120-git://github.com/c9fe/22120.git";
in
{
  nodePackages = super.nodePackages // npmPackages // {
    "a22120" = super.writeShellScriptBin "22120" ''
      ${a22120}/bin/archivist1 $*
    '';
  };
}
