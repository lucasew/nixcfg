{ lib, config, ... }:

{
  config = lib.mkIf config.services.xserver.desktopManager.kodi.enable {
    hardware.opengl.enable = true;
    services = {
      displayManager = {
        sessionData = {
          autologinSession = lib.mkDefault "kodi";
        };
      };
      xserver = {
        enable = true;
        displayManager.lightdm.enable = true;
      };
    };
  };
}
