{ pkgs, global, config, lib, ... }:
let
  inherit (global) username;
in
{

  config = lib.mkIf config.programs.adb.enable {
    users.users.${username}.extraGroups = [ "adbusers" ];
    services.udev.packages = with pkgs; [
      gnome3.gnome-settings-daemon
      android-udev-rules
    ];
  };

}
