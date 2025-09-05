{
  pkgs,
  ...
}:
let
  inherit (pkgs.custom) colors wallpaper;
in
{
  stylix = {
    polarity = if colors.isDark then "dark" else "light";
    image = wallpaper;
    base16Scheme = colors.colors;
    targets = {
      # blender.enable = true;
      # chromium.enable = true;
    };
  };
}
