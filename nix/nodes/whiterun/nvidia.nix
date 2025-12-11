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
    package = config.boot.kernelPackages.nvidiaPackages.stable;
    nvidiaSettings = true;
    open = true;
    # nvidiaPersistenced = true;
  };

  hardware.nvidia-container-toolkit = {
    enable =
      config.virtualisation.docker.enable || config.virtualisation.podman.enable;
    package = pkgs.unstable.nvidia-container-toolkit;
  };


  virtualisation.docker = {
    daemon.settings = {
      runtimes.nvidia = {
        path = "${pkgs.nvidia-container-toolkit.tools}/bin/nvidia-container-runtime";
      };
    };   
  };

  environment.systemPackages = [ pkgs.nvtopPackages.full ];
}
