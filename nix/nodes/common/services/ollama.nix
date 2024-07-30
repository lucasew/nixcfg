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
    services.ollama = {
      acceleration = "cuda";
      listenAddress = "127.0.0.1:${toString config.networking.ports.ollama.port}";
    };

    networking.ports.ollama.enable = true;

    services.nginx.virtualHosts."ollama.${config.networking.hostName}.${config.networking.domain}" = {
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
