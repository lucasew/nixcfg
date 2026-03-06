{
  config,
  pkgs,
  lib,
  ...
}:

lib.mkIf config.programs.sway.enable {
  environment.systemPackages = [ pkgs.swaylock ];
}
