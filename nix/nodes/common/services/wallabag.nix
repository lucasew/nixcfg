{ lib, config, pkgs, ... }:

let
  yaml = pkgs.formats.yaml {};
  parameters = yaml.generate "wallabag-parameters.yaml" {parameters = cfg.config; };

  cfg = config.services.wallabag;
  appDir = pkgs.buildEnv {
    name = "wallabag-app-dir";
    ignoreCollisions = true;
    checkCollisionContents = false;
    paths = [
      (pkgs.runCommand "wallabag-parameters.yaml" {} ''
        mkdir -p $out/config
        ln -sf ${parameters} $out/config/parameters.yml
        substitute ${cfg.package}/app/config/config.yml $out/config/config.yml \
          --replace-fail 'cookie_secure: auto' 'cookie_secure: false' # allow http, I guess
      '')
      "${cfg.package}/app"
    ];
  };

  console = pkgs.writeShellScriptBin "wallabag-console" ''
    if [ "$(whoami)" != "${cfg.user}" ]; then
      exec sudo -u "${cfg.user}" "$0" "$@"
    fi
    export WALLABAG_DATA="${cfg.dataDir}"
    cd "$WALLABAG_DATA"
    exec ${lib.getExe pkgs.php} ${pkgs.wallabag}/bin/console --env=prod $@
  '';
in

{
  options = {
    services.wallabag = {
      enable = lib.mkEnableOption "wallabag";
      config = lib.mkOption {
        description = "Configuration for wallabag";
        type = yaml.type;
        default = {};
      };
      domain = lib.mkOption {
        description = "DNS domain of server";
        type = lib.types.str;
        default = "wallabag.${config.services.ts-proxy.network-domain}";
      };
      user = lib.mkOption {
        description = "Service user";
        type = lib.types.str;
        default = "wallabag";
      };
      group = lib.mkOption {
        description = "Service group";
        type = lib.types.str;
        default = "wallabag";
      };
      dataDir = lib.mkOption {
        description = "Data dir";
        type = lib.types.str;
        default = "/var/lib/wallabag";
      };

      package = lib.mkPackageOption pkgs "wallabag" {};
    };
    
  };

  config = lib.mkIf config.services.wallabag.enable {
    networking.ports.wallabag.enable = true;

    environment.systemPackages = [ console ];
    services.ts-proxy.hosts = {
      wallabag = {
        address = "127.0.0.1:${toString config.networking.ports.wallabag.port}";
        enableTLS = true;
      };
    };
    services.nginx.virtualHosts."${cfg.domain}" = {
      listen = [
        {
          port = config.networking.ports.wallabag.port;
          addr = "127.0.0.1";
        }
      ];
      root = "${cfg.package}/web";
      locations."/" = {
        priority = 10;
        tryFiles = "$uri /app.php$is_args$args";
      };
      locations."~ ^/app\\.php(/|$)" = {
        priority = 100;
        fastcgiParams = {
          SCRIPT_FILENAME = "$realpath_root$fastcgi_script_name";
          DOCUMENT_ROOT = "$realpath_root";
        };
        extraConfig = ''
          fastcgi_pass unix:${config.services.phpfpm.pools.wallabag.socket};
          include ${config.services.nginx.package}/conf/fastcgi_params;
          include ${config.services.nginx.package}/conf/fastcgi.conf;
          internal;
        '';
      };
      locations."~ \\.php$" = {
        priority = 1000;
        return = "404";
      };
    };

    services.redis.servers.wallabag = {
      enable = true;
      inherit (cfg) user;
    };

    services.phpfpm.pools.wallabag = {
      inherit (cfg) user group;
      phpPackage = pkgs.php;
      phpEnv = {
        WALLABAG_DATA = cfg.dataDir;
        PATH = lib.makeBinPath [pkgs.php];
      };
      settings = {
        "listen.owner" = config.services.nginx.user;
        "pm" = "dynamic";
        "pm.max_children" = 32;
        "pm.max_requests" = 500;
        "pm.start_servers" = 1;
        "pm.min_spare_servers" = 1;
        "pm.max_spare_servers" = 5;
        "php_admin_value[error_log]" = "stderr";
        "php_admin_flag[log_errors]" = true;
        "catch_workers_output" = true;
      };
       phpOptions = ''
        extension=${pkgs.phpExtensions.pdo}/lib/php/extensions/pdo.so
        extension=${pkgs.phpExtensions.pdo_pgsql}/lib/php/extensions/pdo_pgsql.so
        extension=${pkgs.phpExtensions.session}/lib/php/extensions/session.so
        extension=${pkgs.phpExtensions.ctype}/lib/php/extensions/ctype.so
        extension=${pkgs.phpExtensions.dom}/lib/php/extensions/dom.so
        extension=${pkgs.phpExtensions.simplexml}/lib/php/extensions/simplexml.so
        extension=${pkgs.phpExtensions.gd}/lib/php/extensions/gd.so
        extension=${pkgs.phpExtensions.mbstring}/lib/php/extensions/mbstring.so
        extension=${pkgs.phpExtensions.xml}/lib/php/extensions/xml.so
        extension=${pkgs.phpExtensions.tidy}/lib/php/extensions/tidy.so
        extension=${pkgs.phpExtensions.iconv}/lib/php/extensions/iconv.so
        extension=${pkgs.phpExtensions.curl}/lib/php/extensions/curl.so
        extension=${pkgs.phpExtensions.gettext}/lib/php/extensions/gettext.so
        extension=${pkgs.phpExtensions.tokenizer}/lib/php/extensions/tokenizer.so
        extension=${pkgs.phpExtensions.bcmath}/lib/php/extensions/bcmath.so
        extension=${pkgs.phpExtensions.intl}/lib/php/extensions/intl.so
        extension=${pkgs.phpExtensions.opcache}/lib/php/extensions/opcache.so
      '';
    };

    services.postgresqlBackup.databases = ["wallabag"];

    services.postgresql = {
      enable = true;
      ensureDatabases = ["wallabag"];
      ensureUsers = [
        {name = "wallabag"; ensureDBOwnership = true;}
      ];
    };

    systemd.services.wallabag-setup = {
      description = "Wallabag setup";
      wantedBy = ["multi-user.target"];
      before = ["phpfpm-wallabag.service"];
      requiredBy = ["phpfpm-wallabag.service"];
      after = ["postgresql.service"];
      path = [pkgs.coreutils pkgs.php pkgs.phpPackages.composer];
      environment = {
        WALLABAG_DATA = cfg.dataDir;
      };

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;
        Type = "oneshot";
        RemainAfterExit = "yes";
        PermissionsStartOnly = true;
      };
      script = ''
      echo "Setting up wallabag files in $WALLABAG_DATA ..."
      cd "${cfg.dataDir}"

      rm -rf var/cache/*
      rm -f app
      ln -sf ${appDir} app
      ln -sf ${cfg.package}/composer.{json,lock} .
      ln -sf ${cfg.package}/{src,translations,templates} .

      if [ ! -f installed ]; then
        echo "Installing wallabag"
        ${lib.getExe console} --env=prod wallabag:install --no-interaction
        touch installed
      else
        ${lib.getExe console} --env=prod doctrine:migrations:migrate --no-interaction
      fi
      ${lib.getExe console} --env=prod cache:clear
    '';
    };

    systemd.tmpfiles.rules = [
      "d ${cfg.dataDir} 0700 ${cfg.user} ${cfg.group} - -"
    ];

    services.rabbitmq.enable = true;
    users.users.${cfg.user} = {
      isSystemUser = true;
      inherit (cfg) group;
    };

    users.groups.${cfg.group} = {};

    services.wallabag.config = {
      database_driver = "pdo_pgsql";
      database_host = null;
      database_port = null;
      database_path = null;
      database_name = "wallabag";
      database_user = cfg.user;
      database_password = null; # socket authenticates using the user
      domain_name = "https://" + cfg.domain;
      database_table_prefix = "wallabag_";
      database_socket = "/run/postgresql/.s.PGSQL.${toString config.services.postgresql.port}";
      database_charset = "utf8";
      server_name = "Wallabag";
      mailer_transport = "sendmail";
      mailer_dsn = "sendmail://default";
      locale = "pt";
      twofactor_auth = false;
      secret = "whatever, it's only available for some nodes in my network anyway";
      fosuser_registration = false;
      fosuser_confirmation = true;
      rss_limit = 50 ;
      rabbitmq_host = "localhost";
      rabbitmq_port = 5672 ;
      rabbitmq_user = "guest" ;
      rabbitmq_password = "guest" ;
      redis_scheme = "tcp" ;
      redis_host = "localhost";
      redis_port = 6379 ;
      redis_path = null ;
      redis_password = null;
      sentry_dsn = null ;
      from_email = "wallabag@${config.networking.domain}";
      twofactor_sender = "wallabag@${config.networking.domain}";
      fos_oauth_server_access_token_lifetime = 3600;
      fos_oauth_server_refresh_token_lifetime = 1209600;
      rabbitmq_prefetch_count = 5;
    };
  };
}
