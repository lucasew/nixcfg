{
  pkgs,
  lib,
  config,
  ...
}:

let
  cfg = config.services.ollama;
in

{

  config = lib.mkIf cfg.enable {
    services.ts-proxy.hosts = {
      ollama = {
        enableTLS = true;
        address = "127.0.0.1:${toString config.networking.ports.ollama-web.port}";
      };
    };

    services.ollama = {
      acceleration = "cuda";
      listenAddress = "127.0.0.1:${toString config.networking.ports.ollama.port}";
    };

    networking.ports.ollama.enable = true;
    networking.ports.ollama-web.enable = true;

    services.nginx.virtualHosts."ollama.${config.networking.hostName}.${config.networking.domain}" = {
      listen = [
        {
          port = config.networking.ports.ollama-web.port;
          addr = "127.0.0.1";
        }
      ];
      root = "${pkgs.ollama-webui}/share/ollama-webui/www";
      locations = {
        "/" = {
          proxyPass = "http://127.0.0.1:${toString config.networking.ports.ollama.port}";
          proxyWebsockets = true;
        };
      };
    };
  };
}
