{ config, ... }: {
  services.transmission = {
    enable = true;
    openFirewall = true;
    openPeerPorts = true;
    settings = {
      # incomplete-dir = "/tmp/transmission/incomplete";
    };
  };
  services.nginx.virtualHosts."transmission.${config.networking.hostName}.${config.networking.domain}" = {
    locations."/" = {
      proxyPass = "http://127.0.0.1:${toString config.services.transmission.settings.rpc-port}";
    };
  };
}
