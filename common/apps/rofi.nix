{ pkgs, config, stdenv, ... }:

let
  # Use the let-in clause to assign the derivation to a variable
  rofiScript = pkgs.writeShellScriptBin "my-rofi" ''
    ${pkgs.rofi}/bin/rofi -show combi -combi-modi window,drun -theme gruvbox-dark -show-icons
  '';
in
rofiScript
