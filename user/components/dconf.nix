{config, pkgs, ...}:
{
    dconf.settings = {
        "org/gnome/desktop/background" = {
            picture-options="zoom";
            picture-uri="file:///nix/store/aslnarw6322gcij310y3cww16sd9fjcf-gnome-backgrounds-3.34.0/share/backgrounds/gnome/Road.jpg";
            primary-color="#ffffff";
            secondary-color="#000000";
        };
        "org/gnome/desktop/input-sources" = {
            current="uint32 0";
            sources=''[("xkb", "br"), ("xkb", "us")]'';
            xkb-options=["terminate:ctrl_alt_bksp"];
        };
        "org/gnome/desktop/interface" = {
            cursor-theme="Adwaita";
            gtk-im-module="gtk-im-context-simple";
            gtk-theme="Adwaita-dark";
        };
        "org/gnome/desktop/peripherals/keyboard" = {
            numlock-state = false;
        };
        "org/gnome/desktop/privacy" = {
            disable-microphone = true;
            report-technical-problems = false;
        };
        "org/gnome/desktop/screensaver" = {
            picture-options="zoom";
            picture-uri="file:///nix/store/aslnarw6322gcij310y3cww16sd9fjcf-gnome-backgrounds-3.34.0/share/backgrounds/gnome/Road.jpg";
            primary-color="#ffffff";
            secondary-color="#000000";
        };
        "org/gnome/system/location" = {
            enabled=false;
        };
    };
}