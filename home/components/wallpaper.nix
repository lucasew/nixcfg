{...}:
let
    globalConfig = import <dotfiles/globalConfig.nix>;
in {
    dconf.settings = {
        "org/gnome/desktop/background" = {
            picture-uri = "file:///${globalConfig.wallpaper}";
        };
    };
    home.file.".background-image".source = globalConfig.wallpaper;
}