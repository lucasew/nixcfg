{ config, lib, ... }:

{
  config = lib.mkIf config.services.restic.server.enable {
    networking.ports.restic-server.enable = true;
    services.restic.server = {
      appendOnly = true;
      listenAddress = "127.0.0.1:${toString config.networking.ports.restic-server.port}";
      extraFlags = [ "--no-auth" ];
    };

    services.ts-proxy.hosts = {
      restic-server = {
        address = "127.0.0.1:${toString config.networking.ports.restic-server.port}";
        enableTLS = true;
      };
    };

    services.prometheus.scrapeConfigs = [
      {
        job_name = "restic-server";
        static_configs = [
          { targets = [ "127.0.0.1:${toString config.networking.ports.restic-server.port}" ]; }
        ];
      }
    ];
  };
}
