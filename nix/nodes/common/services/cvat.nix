{
  lib,
  config,
  pkgs,
  ...
}:

let
  cfg = config.services.cvat;
  inherit (lib) mkEnableOption mkOption types mkIf mkDefault;

  # CVAT image versions
  cvatVersion = "v2.20.1";

  # Helper to create container network name
  networkName = "cvat-network";

in

{
  options.services.cvat = {
    enable = mkEnableOption "CVAT annotation platform";

    user = mkOption {
      description = "Service user";
      type = types.str;
      default = "cvat";
    };

    group = mkOption {
      description = "Service group";
      type = types.str;
      default = "cvat";
    };

    dataDir = mkOption {
      description = "Data directory";
      type = types.str;
      default = "/var/lib/cvat";
    };

    port = mkOption {
      description = "Port for CVAT web interface";
      default = config.networking.ports.cvat.port;
      type = types.port;
    };

    serverImage = mkOption {
      description = "CVAT server image";
      default = "cvat/server:${cvatVersion}";
      type = types.str;
    };

    uiImage = mkOption {
      description = "CVAT UI image";
      default = "cvat/ui:${cvatVersion}";
      type = types.str;
    };

    postgresImage = mkOption {
      description = "PostgreSQL image";
      default = "postgres:15-alpine";
      type = types.str;
    };

    redisImage = mkOption {
      description = "Redis image";
      default = "redis:7.2-alpine";
      type = types.str;
    };
  };

  config = mkIf cfg.enable {
    # Enable port allocation
    networking.ports.cvat.enable = mkDefault true;

    # Create user and group
    users.users.${cfg.user} = {
      isSystemUser = true;
      inherit (cfg) group;
      uid = 2000; # Fixed UID for consistency
    };

    users.groups.${cfg.group} = {
      gid = 2000; # Fixed GID for consistency
    };

    # Create directory structure
    systemd.tmpfiles.rules = [
      "d ${cfg.dataDir} 0755 ${cfg.user} ${cfg.group} - -"
      "d ${cfg.dataDir}/data 0755 ${cfg.user} ${cfg.group} - -"
      "d ${cfg.dataDir}/logs 0755 ${cfg.user} ${cfg.group} - -"
      "d ${cfg.dataDir}/keys 0755 ${cfg.user} ${cfg.group} - -"
      "d ${cfg.dataDir}/postgres 0755 ${cfg.user} ${cfg.group} - -"
      "d ${cfg.dataDir}/redis 0755 ${cfg.user} ${cfg.group} - -"
    ];

    # Create systemd slice with resource limits
    systemd.slices.cvat = {
      description = "CVAT service slice";
      sliceConfig = {
        CPUQuota = "100%"; # 1 VCPU
        MemoryHigh = "1800M";
        MemoryMax = "2048M"; # 2GB RAM
      };
    };

    # Create shared network for containers
    systemd.services."podman-network-${networkName}" = {
      description = "Create podman network for CVAT";
      wantedBy = [ "multi-user.target" ];
      before = [ "cvat.service" ];
      serviceConfig = {
        Type = "oneshot";
        RemainAfterExit = true;
        ExecStart = "${pkgs.podman}/bin/podman network create ${networkName} || true";
        ExecStop = "${pkgs.podman}/bin/podman network rm ${networkName} || true";
      };
    };

    # PostgreSQL container
    virtualisation.oci-containers.containers.cvat-db = {
      image = cfg.postgresImage;
      autoStart = false;
      extraOptions = [
        "--network=${networkName}"
        "--network-alias=cvat-db"
      ];
      volumes = [
        "${cfg.dataDir}/postgres:/var/lib/postgresql/data:Z"
      ];
      environment = {
        POSTGRES_DB = "cvat";
        POSTGRES_USER = "cvat";
        POSTGRES_PASSWORD = "cvat_password";
        POSTGRES_HOST_AUTH_METHOD = "trust";
      };
      user = "2000:2000";
    };

    # Redis container
    virtualisation.oci-containers.containers.cvat-redis = {
      image = cfg.redisImage;
      autoStart = false;
      extraOptions = [
        "--network=${networkName}"
        "--network-alias=cvat-redis"
      ];
      volumes = [
        "${cfg.dataDir}/redis:/data:Z"
      ];
      user = "2000:2000";
    };

    # CVAT Server container
    virtualisation.oci-containers.containers.cvat-server = {
      image = cfg.serverImage;
      autoStart = false;
      dependsOn = [
        "cvat-db"
        "cvat-redis"
      ];
      extraOptions = [
        "--network=${networkName}"
        "--network-alias=cvat-server"
      ];
      volumes = [
        "${cfg.dataDir}/data:/home/django/data:Z"
        "${cfg.dataDir}/logs:/home/django/logs:Z"
        "${cfg.dataDir}/keys:/home/django/keys:Z"
      ];
      environment = {
        DJANGO_MODWSGI_EXTRA_ARGS = "";
        ALLOWED_HOSTS = "*";
        CVAT_REDIS_HOST = "cvat-redis";
        CVAT_REDIS_PORT = "6379";
        CVAT_POSTGRES_HOST = "cvat-db";
        CVAT_POSTGRES_PORT = "5432";
        CVAT_POSTGRES_USER = "cvat";
        CVAT_POSTGRES_DBNAME = "cvat";
        CVAT_POSTGRES_PASSWORD = "cvat_password";
        CVAT_SHARE_URL = "Mounted from /home/django/data directory";
        NUMPROCS = "2";
      };
      user = "2000:2000";
    };

    # CVAT UI container
    virtualisation.oci-containers.containers.cvat-ui = {
      image = cfg.uiImage;
      autoStart = false;
      dependsOn = [ "cvat-server" ];
      extraOptions = [
        "--network=${networkName}"
        "--network-alias=cvat-ui"
      ];
      ports = [
        "127.0.0.1:${toString cfg.port}:80"
      ];
    };

    # CVAT Worker - Low priority
    virtualisation.oci-containers.containers.cvat-worker-low = {
      image = cfg.serverImage;
      autoStart = false;
      dependsOn = [
        "cvat-db"
        "cvat-redis"
      ];
      extraOptions = [
        "--network=${networkName}"
      ];
      volumes = [
        "${cfg.dataDir}/data:/home/django/data:Z"
        "${cfg.dataDir}/logs:/home/django/logs:Z"
        "${cfg.dataDir}/keys:/home/django/keys:Z"
      ];
      environment = {
        CVAT_REDIS_HOST = "cvat-redis";
        CVAT_REDIS_PORT = "6379";
        CVAT_POSTGRES_HOST = "cvat-db";
        CVAT_POSTGRES_PORT = "5432";
        CVAT_POSTGRES_USER = "cvat";
        CVAT_POSTGRES_DBNAME = "cvat";
        CVAT_POSTGRES_PASSWORD = "cvat_password";
        NUMPROCS = "1";
      };
      cmd = [ "run" "worker.low" ];
      user = "2000:2000";
    };

    # CVAT Worker - Default priority
    virtualisation.oci-containers.containers.cvat-worker-default = {
      image = cfg.serverImage;
      autoStart = false;
      dependsOn = [
        "cvat-db"
        "cvat-redis"
      ];
      extraOptions = [
        "--network=${networkName}"
      ];
      volumes = [
        "${cfg.dataDir}/data:/home/django/data:Z"
        "${cfg.dataDir}/logs:/home/django/logs:Z"
        "${cfg.dataDir}/keys:/home/django/keys:Z"
      ];
      environment = {
        CVAT_REDIS_HOST = "cvat-redis";
        CVAT_REDIS_PORT = "6379";
        CVAT_POSTGRES_HOST = "cvat-db";
        CVAT_POSTGRES_PORT = "5432";
        CVAT_POSTGRES_USER = "cvat";
        CVAT_POSTGRES_DBNAME = "cvat";
        CVAT_POSTGRES_PASSWORD = "cvat_password";
        NUMPROCS = "1";
      };
      cmd = [ "run" "worker.default" ];
      user = "2000:2000";
    };

    # Configure systemd services for orchestration
    systemd.services = {
      # Override podman container services to be part of cvat.slice and not autostart
      podman-cvat-db = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [ "podman-network-${networkName}.service" ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      podman-cvat-redis = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [ "podman-network-${networkName}.service" ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      podman-cvat-server = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
        ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      podman-cvat-ui = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [ "podman-cvat-server.service" ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      podman-cvat-worker-low = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
        ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      podman-cvat-worker-default = {
        wantedBy = mkForce [ ];
        partOf = [ "cvat.service" ];
        after = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
        ];
        serviceConfig = {
          Slice = "cvat.slice";
        };
      };

      # Main CVAT service that orchestrates everything
      cvat = {
        description = "CVAT annotation platform";
        wantedBy = mkForce [ ]; # Don't start on boot - start manually when needed
        after = [ "podman-network-${networkName}.service" ];
        requires = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
          "podman-cvat-server.service"
          "podman-cvat-ui.service"
          "podman-cvat-worker-low.service"
          "podman-cvat-worker-default.service"
        ];

        serviceConfig = {
          Type = "oneshot";
          RemainAfterExit = true;
          ExecStart = "${pkgs.coreutils}/bin/echo 'CVAT started'";
          ExecStop = "${pkgs.coreutils}/bin/echo 'CVAT stopped'";
        };
      };

      # Initialization service to run migrations and create superuser
      cvat-init = {
        description = "CVAT initialization (migrations and superuser)";
        after = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
        ];
        requires = [
          "podman-cvat-db.service"
          "podman-cvat-redis.service"
        ];

        path = [ pkgs.podman ];

        serviceConfig = {
          Type = "oneshot";
          RemainAfterExit = true;
          User = cfg.user;
          Group = cfg.group;
        };

        script = ''
          set -e

          # Wait for database to be ready
          echo "Waiting for database to be ready..."
          sleep 10

          # Run migrations
          echo "Running database migrations..."
          podman exec cvat-server python3 manage.py migrate --noinput || true

          # Create superuser if it doesn't exist
          echo "Creating superuser if needed..."
          podman exec cvat-server python3 manage.py shell << 'PYTHON_SCRIPT' || true
from django.contrib.auth import get_user_model
User = get_user_model()
if not User.objects.filter(username='admin').exists():
    User.objects.create_superuser('admin', 'admin@localhost', 'admin')
    print('Superuser created: admin/admin')
else:
    print('Superuser already exists')
PYTHON_SCRIPT

          echo "CVAT initialization complete"
        '';
      };
    };

    # Create cvat-manage script wrapper for manage.py
    environment.systemPackages = [
      (pkgs.writeShellScriptBin "cvat-manage" ''
        # Check if running as cvat user
        if [ "$(whoami)" != "${cfg.user}" ]; then
          exec sudo -u "${cfg.user}" "$0" "$@"
        fi

        # Check if cvat-server container is running
        if ! ${pkgs.podman}/bin/podman ps --format '{{.Names}}' | grep -q '^cvat-server$'; then
          echo "Error: cvat-server container is not running"
          echo "Start CVAT with: systemctl start cvat.service"
          exit 1
        fi

        # Execute manage.py in the container
        exec ${pkgs.podman}/bin/podman exec -it cvat-server python3 manage.py "$@"
      '')
    ];

    # Integrate with ts-proxy
    services.ts-proxy.hosts.cvat = {
      address = "127.0.0.1:${toString cfg.port}";
      enableTLS = true;
      proxies = [ "cvat.service" ];
    };
  };
}
