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
        # enableTLS = true;
        address = "127.0.0.1:${toString config.networking.ports.ollama.port}";
        proxies = [ "ollama.service" ];
      };
    };

    services.ollama = {
      acceleration = "cuda";
      package = pkgs.unstable.ollama;
      environmentVariables = {
        OLLAMA_ORIGINS = "*";
      };
      listenAddress = "0.0.0.0:${toString config.networking.ports.ollama.port}";
    };

    networking.ports.ollama.enable = true;
  };
}
