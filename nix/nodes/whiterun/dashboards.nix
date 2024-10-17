{
  config,
  pkgs,
  lib,
  ...
}:
{

  networking.ports.grafana-web.enable = true;
  # networking.ports.grafana-web.port = lib.mkDefault 49150;
  services.grafana = {
    enable = true;
    settings.server = {
      domain = "grafana.${config.services.ts-proxy.network-domain}";
      http_port = config.networking.ports.grafana-web.port;
      http_addr = "127.0.0.1";
    };
  };

  services.ts-proxy.hosts = {
    grafana = {
      address = "127.0.0.1:${toString config.services.grafana.settings.server.http_port}";
      enableTLS = true;
    };
    prometheus = {
      address = "127.0.0.1:${toString config.services.prometheus.port}";
      enableTLS = true;
    };
  };

  services.nginx.statusPage = true;

  networking.ports.prometheus.enable = true;
  # networking.ports.prometheus.port = lib.mkDefault 49145;
  networking.ports.prometheus-node_exporter.enable = true;
  # networking.ports.prometheus-node_exporter.port = lib.mkDefault 49144;
  services.prometheus = {
    enable = true;
    inherit (config.networking.ports.prometheus) port;

    webExternalUrl = "http://prometheus.${config.services.ts-proxy.network-domain}";
    exporters = {
      node = {
        enable = true;
        enabledCollectors = [
          "logind"
          "systemd"
        ];
        inherit (config.networking.ports.prometheus-node_exporter) port;
      };
      zfs.enable = true;
      nginx = {
        enable = true;
        sslVerify = false; # rede interna não usa ssl
      };
      postgres = {
        enable = true;
        runAsLocalSuperUser = true;
      };
    };
    scrapeConfigs = (
      map
        (item: {
          job_name = item;
          static_configs = [
            { targets = [ "127.0.0.1:${toString config.services.prometheus.exporters.${item}.port}" ]; }
          ];
        })
        [
          "node"
          "zfs"
          "nginx"
          "postgres"
        ]
    );
  };
}
