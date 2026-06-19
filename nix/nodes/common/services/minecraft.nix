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
    systemd.sockets.minecraft-server = {
      wantedBy = lib.mkForce [ ];
    };

    # Inverted dependency (server raises proxy, proxy never raises server):
    # The minecraft-server.service is the thing that activates ts-proxy-mc
    # (and tears it down on stop). This is the opposite of what the `proxies`
    # list in ts-proxy does by default.
    systemd.services.ts-proxy-mc = {
      wantedBy = lib.mkForce [ "minecraft-server.service" ];
      wants = lib.mkForce [ ];
      partOf = lib.mkForce [ "minecraft-server.service" ];
      after = lib.mkForce [ "minecraft-server.service" ];
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
      package = pkgs.unstable.minecraft-server;
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
      # We do not set `proxies` here. Activation/dependency wiring is done
      # explicitly above so that only minecraft-server.service can bring the
      # proxy up (true inversion).
    };
  };
}
