{
  self,
  config,
  pkgs,
  ...
}:
{
  imports = [ "${self.inputs.nixos-hardware}/common/gpu/nvidia" ];

  services.xserver.videoDrivers = [
    "nvidia"
    "modesetting"
  ];

  boot.blacklistedKernelModules = [ "nvidiafb" "nouveau" "nova_core" ]; # CUDA ou GTFO

  hardware.nvidia = {
    package = config.boot.kernelPackages.nvidiaPackages.stable;
    nvidiaSettings = true;
    open = true;
    # nvidiaPersistenced = true;
  };

  hardware.nvidia-container-toolkit = {
    enable =
      config.virtualisation.docker.enable || config.virtualisation.podman.enable;
  };

  environment.systemPackages = [ pkgs.nvtopPackages.full ];
}
