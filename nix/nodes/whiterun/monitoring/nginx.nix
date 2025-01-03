{ config, ... }:

{
  networking.ports.prometheus-exporter-nginx.enable = true;
  services.nginx.statusPage = true;

  services.prometheus = {
    exporters.nginx = {
      enable = true;
      sslVerify = true; # internal net doesn't use ssl
      inherit (config.networking.ports.prometheus-exporter-nginx) port;
    };

    scrapeConfigs = [
      {
        job_name = "nginx";
        static_configs = [
           {
            targets = [
              "127.0.0.1:${toString config.networking.ports.prometheus-exporter-nginx.port}"
            ];
          }
        ];
      }
    ];
  };
}
