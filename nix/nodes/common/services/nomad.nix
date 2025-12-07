{ config, pkgs, lib, ... }:

{
  config = lib.mkIf config.services.nomad.enable {
    services.nomad = {
      extraSettingsPlugins = [ pkgs.nomad-driver-podman ];
      settings = {
        datacenter = lib.mkDefault "local.home.br";
        bind_addr = "{{ GetInterfaceIP \"tailscale0\" }}"; # Dynamically binds to Tailscale IP
        client = {
          enabled = true;
          alloc_mounts_dir = "/var/lib/nomad/alloc_mounts";
          server_join = {
            retry_join = [
              "ravenrock:4647"
            ];
          };
        };
        plugin = [
          {
            nomad-driver-podman = {
              config = {
              };
            };
          }
          {
            nomad-device-nvidia = {
              config = {
                enabled = true;
              };
            };
          }
        ];
      };
    };
  };
}
