{
  config,
  pkgs,
  lib,
  ...
}:

let
  inherit (config.networking.ports.minecraft) port;

in
{
  config = lib.mkIf config.services.minecraft-server.enable {
    networking.ports.minecraft.enable = true;
    systemd.services.minecraft-server = {
      wantedBy = lib.mkForce [ ]; # don't start on boot
    };

    systemd.services.minecraft-server-backup = {
      description = "Minecraft server backup";
      path = [ pkgs.zip ];
      script = ''
        function rcon {
          if [[ -p /run/minecraft-server.stdin ]]; then
            echo "$@" | tee /run/minecraft-server.stdin
          fi
        }
        rcon /say Backup iniciado
        rcon /save-off
        rcon /save-all
        sleep 15

        cd "${config.services.minecraft-server.dataDir}"
        zip -9 -r /var/backup/minecraft.zip world || rcon /say Backup falhou, olha os logs carai

        rcon /save-on
        rcon /say Backup feito
      '';
    };

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
      proxies = [ "minecraft-server.service" ];
      listen = 10000;
    };
  };
}
