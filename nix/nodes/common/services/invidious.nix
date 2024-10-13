{ config, lib, ... }:

lib.mkIf config.services.invidious.enable {
  networking.ports.invidious.enable = true;
  # networking.ports.invidious.port = lib.mkDefault 49149;
  services.invidious = {
    inherit (config.networking.ports.invidious) port;
    settings = {
      db = {
        user = "invidious";
        dbname = "invidious";
      };
    };
  };

  services.ts-proxy.hosts = {
    invidious = {
       addr = "http://127.0.0.1:${toString config.networking.ports.invidious.port}"; 
    };
  };
}
