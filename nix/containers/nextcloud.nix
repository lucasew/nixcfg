{ config, pkgs, lib, ... }:

let
  inherit (lib) mkIf;
  domain = "nextcloud.whiterun.lucao.net"; # Domínio fixo por enquanto
in
{
  config = {
    virtualisation.oci-containers.containers.nextcloud = {
      image = "nextcloud-image";
      ports = [ "8080:80" ];
      volumes = [
        "/var/lib/nextcloud:/var/lib/nextcloud"
        "/var/run/secrets/nextcloud-admin-password:/var/run/secrets/nextcloud-admin-password"
      ];
      environment = {
        POSTGRES_HOST = "postgresql"; # Aponta para o contêiner do postgresql
        POSTGRES_DB = "nextcloud";
        POSTGRES_USER = "nextcloud";
      };

      config = {
        system.stateVersion = "22.11";
        services.nextcloud = {
          package = pkgs.nextcloud28; # Usando uma versão mais recente
          enable = true;
          hostName = domain;
          configureRedis = true;
          config = {
            dbtype = "pgsql";
            dbname = "nextcloud";
            dbuser = "nextcloud";
            dbhost = "postgresql"; # Conecta ao contêiner postgresql
            adminuser = "lucasew";
            adminpassFile = "/var/run/secrets/nextcloud-admin-password";
          };
          settings = {
            trusted_domains = [ domain "localhost" ];
            trusted_proxies = [ "127.0.0.1" "nginx" ]; # Adiciona o proxy nginx
            overwritehost = domain;
            overwriteprotocol = "https";
            "overwrite.cli.url" = "https://''${domain}";
            preview_ffmpeg_path = lib.getExe pkgs.ffmpeg;
            "memories.exiftool" = lib.getExe pkgs.exiftool;
            "memories.ffmpeg_path" = lib.getExe' pkgs.ffmpeg "ffmpeg";
            "memories.ffprobe_path" = lib.getExe' pkgs.ffmpeg "ffprobe";
            "memories.vod.ffmpeg" = lib.getExe' pkgs.ffmpeg "ffmpeg";
            "memories.vod.ffprobe" = lib.getExe' pkgs.ffmpeg "ffprobe";
            recognize = {
              nice_binary = lib.getExe' pkgs.coreutils "nice";
            };
          };
        };

        environment.systemPackages = with pkgs; [
          ffmpeg
          exiftool
          nodejs
          postgresql
        ];

        users.users.nextcloud.extraGroups = [ "render" ];

        systemd.services.nextcloud-setup = {
          script = ''
            ln -sf ''${lib.getExe pkgs.nodejs} ''${config.services.nextcloud.datadir}/store-apps/recognize/bin/node
            ln -sf ''${lib.getExe pkgs.exiftool} ''${config.services.nextcloud.datadir}/store-apps/memories/bin-ext/exiftool/exiftool
          '';
        };
      };
    };
  };
}
