{ pkgs, ... }:
with pkgs.globalConfig;
let
  bin = pkgs.writeShellScriptBin "usb_tixati" ''
    ${pkgs.wine}/bin/wine /run/media/${username}/Dados/PortableApps/PROGRAMAS/Tixati_portable/tixati_Windows32bit.exe
  '';
in {
  home.packages = [bin];
}

