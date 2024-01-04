{ pkgs, lib, config, ... }:

let
  cfg = config.services.ollama;
in

{

  config = lib.mkIf cfg.enable {
    services.ollama.package = pkgs.ollama-cuda;
  
    networking.ports.ollama.enable = true;

    systemd.services.ollama = {
      environment = {
        OLLAMA_HOST = "127.0.0.1:${toString config.networking.ports.ollama.port}";
      };
    };
    
    services.nginx.virtualHosts."ollama.${config.networking.hostName}.${config.networking.domain}" = {
      locations."/" = {
        proxyPass = "http://127.0.0.1:${toString config.networking.ports.ollama.port}";
        proxyWebsockets = true;
      };
    };
  };
}