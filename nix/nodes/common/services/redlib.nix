{
  pkgs,
  lib,
  config,
  ...
}:
let
  cfg = config.services.redlib;
  inherit (lib) mkEnableOption mkOption types mkIf mkDefault;
in
{
  options.services.redlib = {
    enable = mkEnableOption "redlib";
    image = mkOption {
      description = "Which redlib image to use";
      default = "quay.io/redlib/redlib:latest";
      type = types.str;
    };
    port = mkOption {
      description = "Port for redlib";
      default = config.networking.ports.redlib.port;
      type = types.port;
    };
  };

  config = mkIf cfg.enable {
    networking.ports.redlib.enable = mkDefault true;

    services.ts-proxy.hosts = {
      libreddit = {
        address = "127.0.0.1:${toString cfg.port}";
        enableTLS = true;
      };
    };

    virtualisation.oci-containers.containers.redlib = {
      inherit (cfg) image;
      pull = "always";
      ports = [
        "127.0.0.1:${toString cfg.port}:8080"
      ];
      environment = {
        REDLIB_DEFAULT_SHOW_NSFW = "on";
      };
    };
  };
}
