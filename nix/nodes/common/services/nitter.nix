{ config, lib, ... }:

{
  imports = [
    (
      { config, ... }:
      {
        config = lib.mkIf config.services.invidious.enable {
          services.nitter.preferences.replaceYouTube = "invidious.${config.services.ts-proxy.network-domain}";
        };
      }
    )

    (
      { config, ... }:
      {
        config = lib.mkIf config.services.libreddit.enable {
          services.nitter.preferences.replaceReddit = "libreddit.${config.services.ts-proxy.network-domain}";
        };
      }
    )
  ];
  config = lib.mkIf config.services.nitter.enable {
    networking.ports.nitter.enable = true;

    services.nitter.server = {
      inherit (config.networking.ports.nitter) port;
    };

    services.nitter.preferences.replaceTwitter = "nitter.${config.services.ts-proxy.network-domain}";

    services.ts-proxy.hosts = {
      nitter = {
        addr = "http://127.0.0.1:${toString config.networking.ports.nitter.port}";
      };
    };
  };
}
