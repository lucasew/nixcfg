{ lib, config, pkgs, ... }:
let
  inherit (lib) mkIf;
in {
  config = mkIf config.services.nextcloud.enable {
    services.nextcloud.package = pkgs.nextcloud27;
    users.users.nextcloud = {
      extraGroups = [ "admin-password" "render" ];
      packages = with pkgs; [ ffmpeg nodejs coreutils exiftool ];
    };
    services.nextcloud = {
      configureRedis = true;
      hostName = "nextcloud.${config.networking.hostName}.${config.networking.domain}";
      config = {
        dbtype = "pgsql";
        dbname = "nextcloud";
        dbuser = "nextcloud";
        dbhost = "/run/postgresql";
        adminuser = "lucasew";
        adminpassFile = "/var/run/secrets/admin-password";
      };
      extraOptions.enabledPreviewProviders = [
        "OC\\Preview\\BMP"
        "OC\\Preview\\GIF"
        "OC\\Preview\\JPEG"
        "OC\\Preview\\Krita"
        "OC\\Preview\\MarkDown"
        "OC\\Preview\\MP3"
        "OC\\Preview\\OpenDocument"
        "OC\\Preview\\PNG"
        "OC\\Preview\\TXT"
        "OC\\Preview\\XBitmap"
        "OC\\Preview\\HEIC"
      ];
    };

    systemd.services.nextcloud-setup = {
      requires = ["postgresql.service"];
      after = ["postgresql.service"];
    };

    services.postgresqlBackup.databases = [ "nextcloud" ];

    services.postgresql = {
      ensureDatabases = [ "nextcloud" ];
      ensureUsers = [
        {name = "nextcloud"; ensurePermissions."DATABASE nextcloud" = "ALL PRIVILEGES";}
      ];
    };
  };
}
