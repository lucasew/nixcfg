{ config, lib, ... }: {
	config = lib.mkIf config.services.desktopManager.cosmic.enable {
	  services.displayManager.cosmic-greeter.enable = true;
	  services.system76-scheduler.enable = true;
	};
}
