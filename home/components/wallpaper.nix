{...}:
let
    globalConfig = import <dotfiles/globalConfig.nix>;
in {
    dconf.settings = {
        "org/gnome/desktop/background" = {
            picture-uri = "file:///${globalConfig.wallpaper}";
        };
        "org/gnome/desktop/screensaver" = {
          picture-uri = "file:///${globalConfig.wallpaper}";
          picture-options="zoom";
          primary-color="#ffffff";
          secondary-color="#000000";
        };
    };
    home.file.".background-image".source = globalConfig.wallpaper;
}
