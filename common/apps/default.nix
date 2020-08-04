pkgs: rec {
  stremio = pkgs.callPackage ./stremio.nix {};
  my_rofi = pkgs.callPackage ./rofi.nix {};
}
