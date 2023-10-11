{ config, lib, ... }:

{
  config = lib.mkIf config.programs.hyprland.enable {
    programs.regreet.enable = true;
  };
}
