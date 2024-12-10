{
  config,
  pkgs,
  lib,
  ...
}:

let
  cfg = config.services.phputils;
in

{
  options = {
    services.phputils = {
      enable = (lib.mkEnableOption "php teste") // {
        default = true;
      };
      php = lib.mkPackageOption pkgs "php" { };
      scriptDir = lib.mkOption {
        description = "Where are the scripts";
        default = "/etc/phputils";
        type = lib.types.str;
      };
      socket = lib.mkOption {
        description = "Where to listen socket for php test";
        default = "/run/php-test.sock";
      };
    };
  };
  config = lib.mkIf cfg.enable {

    systemd.sockets.phputils = {
      restartTriggers = [ cfg.socket ];
      socketConfig = {
        ListenStream = cfg.socket;
        Accept = true;
      };
      partOf = [ "phputils.service" ];
      wantedBy = [
        "sockets.target"
        "multi-user.target"
      ];
    };

    systemd.slices.phputils.sliceConfig = {
      MemoryMax = "64M";
      MemoryHigh = "16M";
      CPUQuota = "10%";
      ManagedOOMSwap = "kill";
      ManagedOOMPressure = "kill";
    };

    systemd.services."phputils@" = {
      stopIfChanged = true;
      after = [ "network.target" ];
      serviceConfig = {
        Slice = "phputils.slice";
        StandardInput = "socket";
        StandardOutput = "socket";
        StandardError = "journal";

        DevicePolicy = "closed";
        MemoryDenyWriteExecute = true;
        NoNewPrivileges = true;
        PrivateDevices = true;
        PrivateTmp = true;
        ProtectControlGroups = true;
        # ProtectHome = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectKernelLogs = true;
        ProtectSystem = "strict";
      };

      script = ''
        cd "${cfg.scriptDir}"
        export SCRIPT_DIR="${cfg.scriptDir}"
        exec ${lib.getExe cfg.php}  -d display_errors="stderr" -d disable_functions="header" ${./entrypoint.php}
      '';
    };

    systemd.tmpfiles.rules = [ "d ${cfg.scriptDir} 0700 root root - -" ];

    services.ts-proxy.hosts."phputils-${config.networking.hostName}" = {
      enableTLS = true;
      enableFunnel = true;
      network = "unix";
      address = cfg.socket;
    };
  };
}
