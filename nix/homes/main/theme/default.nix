{
  pkgs,
  ...
}:
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
}
