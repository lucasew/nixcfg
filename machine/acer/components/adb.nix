{config, pkgs, ...}:

let 
    common = import ./common;
in 
{
    programs.adb.enable = true;
    users.users.${common.username}.extraGroups = ["adbusers"];
    services.udev.packages = [
        pkgs.android-udev-rules
    ];
}