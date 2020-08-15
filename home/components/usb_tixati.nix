{ pkgs ? import <nixpkgs> {}
, ... }:
let
  cfg = import ../../config.nix;
  bin = pkgs.writeShellScriptBin "usb_tixati" ''
    ${pkgs.wine}/bin/wine /run/media/${cfg "username"}/Dados/PortableApps/PROGRAMAS/Tixati_portable/tixati_Windows32bit.exe
  '';
in {
  home.packages = [bin];
}

