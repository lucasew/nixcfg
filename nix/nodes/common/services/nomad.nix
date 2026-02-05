{
  config,
  pkgs,
  lib,
  ...
}: {
  config = lib.mkIf config.services.nomad.enable {
    services.nomad = {
      extraSettingsPlugins = [
        pkgs.nomad-driver-podman
        pkgs.nomad-driver-nvidia
      ];
      dropPrivileges = false;
      settings = {
        data_dir = "/var/lib/private/nomad";
        datacenter = lib.mkDefault "br_home_local";
        bind_addr = "{{ GetInterfaceIP \"tailscale0\" }}"; # Dynamically binds to Tailscale IP
        "advertise" = {
          "http" = "{{ GetInterfaceIP \"tailscale0\" }}";
          "rpc" = "{{ GetInterfaceIP \"tailscale0\" }}";
          "serf" = "{{ GetInterfaceIP \"tailscale0\" }}";
        };
        client = {
          enabled = true;
          network_interface = "tailscale0";
          host_volume."consul-data" = {
            path = "/var/lib/nomad/consul";
            read_only = false;
          };
          alloc_mounts_dir = "/var/lib/private/nomad/alloc_mounts";
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
          {
            docker = {
              config = {
                image_pull_timeout = "30m";
              };
            };
          }
        ];
      };
    };
  };
}
