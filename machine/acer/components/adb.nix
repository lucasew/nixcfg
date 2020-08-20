{ config, ... }:
let
  pkgs = import <nixpkgs> {};
in
{
  programs.adb.enable = true;
  services.udev.packages = [
    pkgs.android-udev-rules
  ];
  users.users.${pkgs.globalConfig.username}.extraGroups = [ "adbusers" ];
}
