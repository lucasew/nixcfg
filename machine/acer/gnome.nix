{ config, pkgs, ...}:
{
	services.xserver.displayManager.gdm.enable = true;
	services.xserver.desktopManager.gnome3.enable = true;
}
