{
  pkgs,
  lib,
  config,
  ...
}:
lib.mkIf config.services.redlib.enable {
  networking.ports.redlib.enable = true;

  services.redlib = {
    package = pkgs.unstable.redlib;
    inherit (config.networking.ports.redlib) port;
  };

  services.ts-proxy.hosts = {
    libreddit = {
      address = "127.0.0.1:${toString config.services.redlib.port}";
      enableTLS = true;
    };
  };
}
