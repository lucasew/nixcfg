{ pkgs, config, stdenv, ... }:
pkgs.writeShellScriptBin "usb_tixati" ''
${pkgs.wine}/bin/wine /run/media/lucasew/Dados/PortableApps/PROGRAMAS/Tixati_portable/tixati_Windows32bit.exe
''