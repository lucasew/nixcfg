{pkgs, config, ... }:
let
in
{
  services.xserver.displayManager.gdm.enable = true;
  services.xserver.desktopManager.gnome3.enable = true;
  # xdg.portal= {
  #   enable = true;
  #   gtkUsePortal = true;
  # };
}
