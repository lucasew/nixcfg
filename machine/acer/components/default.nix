{ config, pkgs, ... }: {
  imports = [
    ./gui
    ./adb.nix
    ./anbox.nix
  ];
}
