let
  pkgs = import <nixpkgs> {};
in
name:
let
  pkg = pkgs."${name}";
in "${pkg}/bin/${name}"
