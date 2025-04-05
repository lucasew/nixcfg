{
  self,
  config,
  pkgs,
  ...
}:
{
  imports = [ "${self.inputs.nixos-hardware}/common/gpu/nvidia" ];

  services.xserver.videoDrivers = [
    "modesetting"
    "nvidia"
  ];

  hardware.nvidia = {
    package = config.boot.kernelPackages.nvidiaPackages.production; # slightly older, doesn't require to be the open
    nvidiaSettings = true;
    # open = true;
    # nvidiaPersistenced = true;
  };

  hardware.nvidia-container-toolkit.enable =
    config.virtualisation.docker.enable || config.virtualisation.podman.enable;

  # boot.initrd.kernelModules = [ "nvidia" ];
  # boot.extraModulePackages = [ config.boot.kernelPackages.nvidia_x11 ];

  environment.systemPackages = [ pkgs.nvtopPackages.full ];
}
