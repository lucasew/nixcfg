{
  pkgs,
  lib,
  config,
  ...
}:

{
  options.programs.discord-custom.enable = lib.mkEnableOption "discord-desktop";

  config = lib.mkIf config.programs.discord-custom.enable { home.packages = [ pkgs.discord ]; };
}
