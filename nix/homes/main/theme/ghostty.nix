{ pkgs, lib, ... }:

let
  inherit (pkgs.custom) colors;
  inherit (builtins.mapAttrs (k: v: lib.toLower v) colors.colors)
    base00
    base01
    base02
    base03
    base04
    base05
    base06
    base07
    base08
    base09
    base0A
    base0B
    base0C
    base0D
    base0E
    base0F
    ;
in

{
  programs.ghostty.themes.base16-custom = {
    

    background = "#${base00}";
    foreground = "#${base05}";

    selection-background = "#${base02}";
    selection-foreground = "#${base00}";

    palette = [
      "0=#${base00}"
      "1=#${base08}"
      "2=#${base0B}"
      "3=#${base0A}"
      "4=#${base0D}"
      "5=#${base0E}"
      "6=#${base0C}"
      "7=#${base05}"
      "8=#${base03}"
      "9=#${base08}"
      "10=#${base0B}"
      "11=#${base0A}"
      "12=#${base0D}"
      "13=#${base0E}"
      "14=#${base0C}"
      "15=#${base07}"
      "16=#${base09}"
      "17=#${base0F}"
      "18=#${base01}"
      "19=#${base02}"
      "20=#${base04}"
      "21=#${base06}"
    ];
  };
}
