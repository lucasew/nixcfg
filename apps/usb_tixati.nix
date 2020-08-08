{ pkgs, config, stdenv, ... }:
let
  cfg = import ../config.nix;
in
pkgs.writeShellScriptBin "usb_tixati" ''
  ${pkgs.wine}/bin/wine /run/media/${cfg "username"}/Dados/PortableApps/PROGRAMAS/Tixati_portable/tixati_Windows32bit.exe
''
