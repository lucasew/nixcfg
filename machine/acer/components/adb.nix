{ config, ... }:
let
  pkgs = import <dotfiles/pkgs.nix>;
in
{
  programs.adb.enable = true;
  services.udev.packages = [
    pkgs.android-udev-rules
  ];
  users.users.${pkgs.globalConfig.username}.extraGroups = [ "adbusers" ];
}
