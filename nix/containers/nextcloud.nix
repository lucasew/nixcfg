{ config, pkgs, lib, ... }:

let
  domain = "nextcloud.whiterun.lucao.net"; # Domínio fixo por enquanto

  # 1. Definir a configuração NixOS para a imagem
  nextcloud-nixos-config = {
    system.stateVersion = "22.11";
    services.nextcloud = {
      package = pkgs.nextcloud28;
      enable = true;
      hostName = domain;
      configureRedis = true;
      config = {
        dbtype = "pgsql";
        dbname = "nextcloud";
        dbuser = "nextcloud";
        dbhost = "postgresql"; # Conecta-se ao contêiner postgresql pelo nome do host
        adminuser = "lucasew";
        adminpassFile = "/run/secrets/nextcloud-admin-password"; # Caminho dentro do contêiner
      };
      settings = {
        trusted_domains = [ domain "localhost" ];
        trusted_proxies = [ "127.0.0.1" ];
        overwritehost = domain;
        overwriteprotocol = "https";
        "overwrite.cli.url" = "https://${domain}";
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
      # O cliente psql pode ser útil para depuração dentro do contêiner
      postgresql
    ];

    users.users.nextcloud.extraGroups = [ "render" ];

    # Este script de setup ainda é útil dentro da imagem
    systemd.services.nextcloud-setup = {
      script = ''
        ln -sf ${lib.getExe pkgs.nodejs} ${config.services.nextcloud.datadir}/store-apps/recognize/bin/node
        ln -sf ${lib.getExe pkgs.exiftool} ${config.services.nextcloud.datadir}/store-apps/memories/bin-ext/exiftool/exiftool
      '';
    };

    # A imagem precisa de uma rede funcional
    networking.hostName = "nextcloud";
  };

  # 2. Construir a imagem OCI
  nextcloud-image = pkgs.dockerTools.buildImage {
    name = "nixos-nextcloud";
    tag = "latest";
    config = {
      Cmd = [ "${pkgs.nixos-container}/bin/nixos-container" "run-unconfigured-container" ];
    };
    nixOSConfiguration = nextcloud-nixos-config;
  };

in
{
  # 3. Definir o contêiner para usar a imagem construída
  virtualisation.oci-containers.containers.nextcloud = {
    imageFile = nextcloud-image;
    ports = [ "8080:80" ]; # Mapeia a porta 80 do contêiner para 8080 no host
    volumes = [
      "/var/lib/nextcloud:/var/lib/nextcloud"
      "/var/run/secrets/nextcloud-admin-password:/run/secrets/nextcloud-admin-password:ro"
    ];
    environment = {
      # As variáveis de ambiente são passadas para o contêiner, mas a configuração
      # do Nextcloud já está definida para usar o host 'postgresql'.
      # Manter isso pode ser útil para outros scripts ou depuração.
      POSTGRES_HOST = "postgresql";
      POSTGRES_DB = "nextcloud";
      POSTGRES_USER = "nextcloud";
    };
  };
}
