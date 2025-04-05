{ config, ... }:

{
  networking.ports.prometheus-exporter-postgres.enable = true;

  services.prometheus = {
    exporters.postgres = {
      enable = true;
      runAsLocalSuperUser = true;
      inherit (config.networking.ports.prometheus-exporter-postgres) port;
    };

    scrapeConfigs = [
      {
        job_name = "postgres";
        static_configs = [
          {
            targets = [
              "127.0.0.1:${toString config.networking.ports.prometheus-exporter-postgres.port}"
            ];
          }
        ];
      }
    ];
  };
}
