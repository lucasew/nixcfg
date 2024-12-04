{lib, config, pkgs, ...}:

let
  cfg = config.services.ts-proxy;
in

{
  options = {
    services.ts-proxy = {
      network-domain = lib.mkOption {
        description = "Which ts.net domain this machine belongs";
        default = "stargazer-shark.ts.net";
      };

      environmentFile = lib.mkOption {
        description = "Path to environment file for ts-proxy credentials";
        default = "/run/secrets/ts-proxy";
      };

      package = lib.mkPackageOption pkgs "ts-proxy" {};

      user = lib.mkOption {
        description = "Service user";
        type = lib.types.str;
        default = "tsproxy";
      };
      group = lib.mkOption {
        description = "Service group";
        type = lib.types.str;
        default = "tsproxy";
      };

      dataDir = lib.mkOption {
        description = "Data dir";
        type = lib.types.str;
        default = "/var/lib/ts-proxy";
      };

      hosts = lib.mkOption {
        description = "Services to expose to ts-proxy";

        type = lib.types.attrsOf (lib.types.submodule ({ name, ...}: {
          options = {
            enableFunnel = lib.mkEnableOption "enable funnel for this endpoint";
            enableTLS = lib.mkEnableOption "enable TLS for this endpoint";
            enableRaw = lib.mkEnableOption "treat this endpoint as a raw TCP socket";

            network = lib.mkOption {
              description = "First parameter of net.Dial";
              type = lib.types.str;
              default = "";
            };

            address = lib.mkOption {
              description = "Second parameter of net.Dial";
              type = lib.types.str;
            };

            listen = lib.mkOption {
              description = "Which port to listen in the vhost";
              type = lib.types.port;
              default = 0;
            };

            name = lib.mkOption {
              description = "Service name";
              type = lib.types.str;
              default = name;
            };

            unitName = lib.mkOption {
              description = "Systemd unit of the proxy";
              type = lib.types.str;
              default = "ts-proxy-${name}";
            };
          };
          
        }));
      };
    };
  };

  config = {
    sops.secrets.ts-proxy = {
        sopsFile = ../../../../secrets/ts-proxy.env;
        owner = cfg.user;
        group = cfg.group;
        format = "dotenv";
      };

      users.users.${cfg.user} = {
        isSystemUser = true;
        inherit (cfg) group;
      };

      users.groups.${cfg.group} = {};

      systemd.tmpfiles.rules = [
        "d ${cfg.dataDir} 0700 ${cfg.user} ${cfg.group} - -"
      ];

    systemd.slices.ts-proxys.sliceConfig = {
      CPUQuota = "10%";
      MemoryHigh = "256M";
      MemoryMax = "384M";
    };

    systemd.services = lib.mkMerge (builtins.attrValues (builtins.mapAttrs (k: host: {      
      ${host.unitName} = {
        
        description = "ts-proxy service for ${host.name}";
        wantedBy = ["multi-user.target"];

        restartIfChanged = true;

        serviceConfig = {
          Slice = "ts-proxy.slice";
          User = cfg.user;
          Group = cfg.group;
          Restart = "always";
          RestartSec = "10s";
          EnvironmentFile = cfg.environmentFile;
        };

        script = ''
          ${lib.getExe' pkgs.ts-proxy "ts-proxyd"} ${lib.escapeShellArgs ([]
            ++ (["-address" host.address])
            ++ (lib.optional host.enableFunnel "-f")
            ++ (lib.optionals (host.listen != 0) ["-listen" ":${toString host.listen}"])
            ++ (["-n" host.name])
            ++ (lib.optionals (host.network != "") ["-net" host.network])
            ++ (lib.optional host.enableRaw "-raw")
            ++ (["-s" "${cfg.dataDir}/tsproxy-${host.name}"])
            ++ (lib.optional host.enableTLS "-t")
          )}
        '';
      };
    }) cfg.hosts));
  };
}
