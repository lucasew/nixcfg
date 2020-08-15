{ config, pkgs, ... }: {
  imports = [
    ./adb.nix
    ./anbox.nix
    ./gui
  ];
}
