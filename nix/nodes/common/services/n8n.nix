{ config, lib, ...}:

{
  config = lib.mkIf config.services.n8n.enable {
    networking.ports.n8n.enable = true;

    services.ts-proxy.hosts.n8n = {
      enableTLS = true;
      enableFunnel = true;
      address = "localhost:${toString config.services.n8n.settings.port}";
      proxies = [ "n8n.service" ];
    };
    services.n8n = {
      settings = {
        inherit (config.networking.ports.n8n) port;
      };
      webhookUrl = "https://n8n.${config.services.ts-proxy.network-domain}";
    };
    systemd.services.n8n = {
      environment = {
        N8N_PORT = toString config.services.n8n.settings.port;
      };
      serviceConfig = {
        User = "n8n";
        Group = "n8n";
      };
    };
  };
}
