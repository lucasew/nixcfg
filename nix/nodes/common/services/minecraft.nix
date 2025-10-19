{
  config,
  pkgs,
  lib,
  ...
}:

let
  cfg = config.services.minecraft-server;
  inherit (config.networking.ports.minecraft) port;
  dataDir = "/var/lib/minecraft";

in
{
  config = lib.mkIf cfg.enable {
    networking.ports.minecraft.enable = true;

    # OCI container for Minecraft server
    virtualisation.oci-containers.containers.minecraft = {
      image = "itzg/minecraft-server:java8";
      autoStart = false; # don't start on boot - only when manually started
      volumes = [
        "${dataDir}:/data"
      ];
      environment = {
        EULA = "TRUE";
        VERSION = "1.15.2";
        TYPE = "VANILLA";
        DIFFICULTY = "hard";
        MODE = "survival";
        MAX_PLAYERS = "10";
        MOTD = "lucaocraft";
        ONLINE_MODE = "FALSE";
        SERVER_PORT = toString port;
      };
      ports = [
        "127.0.0.1:${toString port}:${toString port}"
      ];
      extraOptions = [
        "--pull=always" # always pull latest image
      ];
    };

    # Ensure data directory exists
    systemd.tmpfiles.rules = [
      "d ${dataDir} 0755 root root - -"
    ];

    # Backup service for Minecraft
    systemd.services.minecraft-backup = {
      description = "Minecraft server backup";
      path = [ pkgs.zip pkgs.podman ];
      script = ''
        function rcon {
          # Send RCON commands to the container
          ${pkgs.podman}/bin/podman exec minecraft rcon-cli "$@" 2>/dev/null || true
        }

        rcon say Backup iniciado
        rcon save-off
        rcon save-all flush
        sleep 15

        cd "${dataDir}"
        zip -9 -r /var/backup/minecraft.zip world || rcon say Backup falhou, olha os logs carai

        rcon save-on
        rcon say Backup feito
      '';
    };

    # Configure ts-proxy to start/stop with Minecraft container
    services.ts-proxy.hosts.mc = {
      enableRaw = true;
      enableFunnel = true;
      address = "127.0.0.1:${toString port}";
      proxies = [ "podman-minecraft.service" ]; # link to container service
      listen = 10000;
    };
  };
}
