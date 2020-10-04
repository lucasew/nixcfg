{pkgs, config, ... }:
let
in
{
  services.xserver.displayManager.gdm.enable = true;
  services.xserver.desktopManager.gnome3.enable = true;
  environment.systemPackages = with pkgs; [
    gnomeExtensions.night-theme-switcher
    gnomeExtensions.sound-output-device-chooser
    gnomeExtensions.gsconnect
  ];
  # xdg.portal= {
  #   enable = true;
  #   gtkUsePortal = true;
  # };
}
