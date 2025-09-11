{
  pkgs,
  self,
  lib,
  config,
  ...
}:
let
  cfg = config.services.cf-torrent;
  inherit (lib)
    mkEnableOption
    mkOption
    types
    mkIf
    mkForce
    mkDefault
    ;
in
{
  options.services.cf-torrent = {
    enable = mkEnableOption "cf-torrent";
    image = mkOption {
      description = "Which cf-torrent image to use";
      default = "ghcr.io/lucasew/cf-torrent:latest";
      type = types.str;
    };
    port = mkOption {
      description = "Port for cf-torrent";
      default = config.networking.ports.cf-torrent.port;
      type = types.port;
    };
    shutdownTimeout = mkOption {
      description = "Time in ms to shutdown the service when inactive";
      default = 10;
      type = types.int;
    };
  };

  config = mkIf cfg.enable {
    networking.ports.cf-torrent.enable = mkDefault true;

    services.cf-torrent.port = mkDefault config.networking.ports.cf-torrent.port;

    services.ts-proxy.hosts = {
      cf-torrent = {
        address = "127.0.0.1:${toString cfg.port}";
        enableTLS = true;
        proxies = [ "cf-torrent.socket" ];
      };
    };

    systemd.sockets.cf-torrent = {
      socketConfig = {
        ListenStream = cfg.port;
      };
      partOf = [
        "cf-torrent.service"
      ];
      wantedBy = [
        "sockets.target"
        "multi-user.target"
      ];
    };

    virtualisation.oci-containers.containers.cf-torrent = {
      inherit (cfg) image;
      environment = {
        IDLE_TIMEOUT = toString cfg.shutdownTimeout;
      };
      pull = "always";
      serviceName = "cf-torrent";
      autoStart = false;
    };

    systemd.services.cf-torrent.serviceConfig.Restart = mkForce "no";

  };
}
