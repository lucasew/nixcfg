{
  config,
  ...
}:
{

  imports = [
    ./nginx.nix
    ./node-exporter.nix
    ./zfs.nix
    ./postgres.nix
    ./nvidia.nix
  ];

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

  networking.ports.prometheus.enable = true;
  services.prometheus = {
    enable = true;
    inherit (config.networking.ports.prometheus) port;

    webExternalUrl = "http://prometheus.${config.services.ts-proxy.network-domain}";
  };
}
