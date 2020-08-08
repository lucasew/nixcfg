pkgs: rec {
  stremio = pkgs.callPackage ./stremio.nix {};
  my_rofi = pkgs.callPackage ./rofi.nix {};
  usb_tixati = pkgs.callPackage ./usb_tixati.nix {};
}
