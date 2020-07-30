{pkgs, ...}:

{
    allowUnfree = true;
    packageOverrides = pkgs: rec {
        stremio = pkgs.callPackage ./common/apps/stremio.nix {};
    };
}
