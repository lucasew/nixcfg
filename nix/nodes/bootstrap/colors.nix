{
  pkgs,
  lib,
  self,
  ...
}:
let
  inherit (self) colors;
in
{
  stylix = {
    polarity = if colors.isDark then "dark" else "light";
    base16Scheme = colors.colors;
  };
}
