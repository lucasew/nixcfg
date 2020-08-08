{ config, pkgs, ... }:

let
  cfg = import ../../../config.nix;
in
{
  programs.adb.enable = true;
  users.users.${cfg "username"}.extraGroups = [ "adbusers" ];
  services.udev.packages = [
    pkgs.android-udev-rules
  ];
}
