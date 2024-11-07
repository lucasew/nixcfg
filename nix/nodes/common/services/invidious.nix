{ config, lib, pkgs, ... }:

lib.mkIf config.services.invidious.enable {
  networking.ports.invidious.enable = true;
  # networking.ports.invidious.port = lib.mkDefault 49149;
  systemd.services.invidious.serviceConfig.DynamicUser = lib.mkForce false;

  sops.secrets."invidious" = {
    sopsFile = ../../../secrets/invidious.txt;
    owner = "invidious";
    group = "invidious";
    format = "binary";
  };
  services.invidious = {
    package = pkgs.unstable.invidious;
    extraSettingsFile = "/var/run/secrets/invidious";
    inherit (config.networking.ports.invidious) port;
    settings = {
      db = {
        user = "invidious";
        dbname = "invidious";
      };
    };
  };

  users = {
    users.invidious = {
      isSystemUser = true;
      group = "invidious";
    };
    groups.invidious = {};
  };

  services.ts-proxy.hosts = {
    invidious = {
      address = "127.0.0.1:${toString config.networking.ports.invidious.port}"; 
      enableTLS = true;
    };
  };
}
