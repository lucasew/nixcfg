{ config, pkgs, lib, ... }:

{
  config = lib.mkIf config.services.nomad.enable {
    services.nomad = {
      extraSettingsPlugins = [
        pkgs.nomad-driver-podman
        pkgs.nomad-driver-nvidia
      ];
      settings = {
        datacenter = lib.mkDefault "br_home_local";
        bind_addr = "{{ GetInterfaceIP \"tailscale0\" }}"; # Dynamically binds to Tailscale IP
        client = {
          enabled = true;
          host_volume."consul-data" = {
            path      = "/var/lib/nomad/consul";
            read_only = false;
          };
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
