{ config, lib, ... }:

{

  imports = [
    (
      { config, ... }:
      {
        config = lib.mkIf config.services.invidious.enable {
          services.miniflux.config = {
            INVIDIOUS_INSTANCE = "https://invidious.${config.services.ts-proxy.network-domain}";
            YOUTUBE_EMBED_URL_OVERRIDE = "https://invidious.${config.services.ts-proxy.network-domain}/embed/";
          };
        };
      }
    )
  ];
  config = lib.mkIf config.services.miniflux.enable {
    networking.ports.miniflux.enable = true;

    services.miniflux = {
      config = {
        LISTEN_ADDR = "localhost:${toString config.networking.ports.miniflux.port}";
        BASE_URL = "https://miniflux.${config.services.ts-proxy.network-domain}";
        FETCH_ODYSEE_WATCH_TIME = toString 1;
        FETCH_YOUTUBE_WATCH_TIME = toString 1;
      };

      # if you are not allowed you shouldn't even been able to open the homepage lol
      adminCredentialsFile = builtins.toFile "creds" ''
        ADMIN_USERNAME=admin
        ADMIN_PASSWORD=adminadmin
      '';
    };

    services.ts-proxy.hosts = {
      miniflux = {
        enableTLS = true;
        address = "127.0.0.1:${toString config.networking.ports.miniflux.port}";
      };
    };
    services.postgresqlBackup.databases = [ "miniflux" ];
  };
}
