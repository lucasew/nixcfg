{ config, ... }:
let
  pkgs = import <dotfiles/pkgs.nix>;
in
{
  services.xserver.displayManager.gdm.enable = true;
  services.xserver.desktopManager.gnome3.enable = true;
  xdg.portal.enable = true;
  # xdg.portal.gtkUsePortal = true;
}
