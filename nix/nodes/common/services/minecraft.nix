{ config, pkgs, lib, ... }:

let
  inherit (config.networking.ports.minecraft) port;

in
{
  config = lib.mkIf config.services.minecraft-server.enable {
    networking.ports.minecraft.enable = true;

    services.minecraft-server = {
      package = pkgs.unstable.minecraftServers.vanilla-1-15.override {
        jre_headless = pkgs.unstable.openjdk8;
      };
      declarative = true;
      eula = true;
      serverProperties = {
        server-port = port;
        difficulty = 3;
        gamemode = 0; # survival
        max-players = 10;
        motd = "lucaocraft";
        online-mode = false;
      };
    };
    services.ts-proxy.hosts.mc = {
      enableRaw = true;
      enableFunnel = true;
      address = "127.0.0.1:${toString port}";
      listen = 10000;
    };
  };
}
