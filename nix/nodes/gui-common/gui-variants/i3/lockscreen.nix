{
  pkgs,
  config,
  lib,
  ...
}:

lib.mkIf config.services.xserver.windowManager.i3.enable {
  environment.systemPackages = [ pkgs.xss-lock ];
}
