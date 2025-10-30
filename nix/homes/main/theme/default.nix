{
  pkgs,
  config,
  ...
}:
let
  inherit (pkgs.custom.colors-lib-contrib) gtkThemeFromScheme shellThemeFromScheme;
  inherit (pkgs.custom) colors colorpipe;
in
{
  imports = [
    # ./terminator.nix
    # ./telegram
    # ./obsidian
    # ./discord
    # ./qtct.nix
    # ./ghostty.nix
    ../../../stylix.nix
  ];

  home.packages = with pkgs; [
    lxappearance
    colorpipe
  ];

  gtk = {
    # theme = {
    #   package = gtkThemeFromScheme { scheme = colors; };
    #   name = colors.slug;
    # };
    # cursorTheme = {
    #   package = pkgs.paper-icon-theme;
    #   name = "Paper";
    # };
    # iconTheme = {
    #   package = pkgs.paper-icon-theme;
    #   name = "Paper";
    # };
  };
  programs.bash.bashrcExtra = ''
    function setup_colors {
      ${shellThemeFromScheme { scheme = config.stylix.generated // {slug = "stylix";}; }}
    }
    setup_colors
  '';
}
