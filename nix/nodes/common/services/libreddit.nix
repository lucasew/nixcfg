{ pkgs, lib, config, ... }:
lib.mkIf config.services.libreddit.enable {
  networking.ports.libreddit.enable = true;
  # networking.ports.libreddit.port = lib.mkDefault 49147;

  services.libreddit = {
    package = pkgs.unstable.redlib;
    inherit (config.networking.ports.libreddit) port;
  };

  services.ts-proxy.hosts = {
    libreddit = {
      address =  "127.0.0.1:${toString config.services.libreddit.port}";
      enableTLS = true;
    };
  };
}
