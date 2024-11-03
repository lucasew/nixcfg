{
  config,
  pkgs,
  lib,
  ...
}:

let
  inherit (lib)
    mkEnableOption
    mkIf
    mkOption
    types
    ;

  cfg = config.services.rsyncnet-remote-backup;

  wrappedSsh = pkgs.writeShellScriptBin "wssh" ''
    exec ${lib.getExe pkgs.openssh}  \
      -v \
      -T \
      -o StrictHostKeyChecking=no \
      -o IdentitiesOnly=yes \
      -i /run/secrets/rsyncnet-remote-backup \
      "${cfg.host}" \
      -- \
      "$@"
  '';
in

{
  options.services.rsyncnet-remote-backup = {
    enable = mkEnableOption "rsync.net remote backup";
    user = mkOption {
      default = "rsyncnet";
      description = "User for the rsyncnet stuff";
      type = types.str;
    };
    group = mkOption {
      default = "rsyncnet";
      description = "Group for the rsyncnet stuff";
      type = types.str;
    };
    host = mkOption {
      default = "de3163@de3163.rsync.net";
      description = "Which rsync.net account/user";
      type = types.str;
    };
    git-step-timeout = mkOption {
      default = 600;
      description = "Timeout to run a git job";
      type = types.int;
    };
    calendar = mkOption {
      default = "00:00:01";
      description = "When to run the backups";
      type = types.str;
    };
    dataDir = lib.mkOption {
      description = "Data dir";
      type = lib.types.str;
      default = "/var/lib/rsyncnet-items";
    };
  };

  config = mkIf cfg.enable {
    users = {
      users.${cfg.user} = {
        isSystemUser = true;
        group = cfg.group;
        home = cfg.dataDir;
      };
      groups.${cfg.group} = { };
    };

    sops.secrets.rsyncnet-remote-backup = {
      sopsFile = ../../../../secrets/rsyncnet;
      owner = cfg.user;
      group = cfg.group;
      format = "binary";
    };

    systemd.tmpfiles.rules = [
      "d ${cfg.dataDir} 0700 ${cfg.user} ${cfg.group} - -"
    ];

    systemd.timers."rsyncnet-remote-backup" = {
      description = "rsync.net backup timer";
      wantedBy = [ "timers.target" ];
      timerConfig = {
        OnCalendar = cfg.calendar;
        AccuracySec = "30m";
        Unit = "rsyncnet-remote-backup.service";
      };
    };

    systemd.services."rsyncnet-remote-backup" = {
      path = with pkgs; [ wrappedSsh bash pv git gawk rsync openssh ];

      restartIfChanged = false;
      stopIfChanged = false;

      script = ''
        cd "${cfg.dataDir}"
        export PATH+=":/run/wrappers/bin"


        for repo in $(wssh ls git); do
          unit="rsyncnet-remote-backup-git@$repo"
          systemctl start "$unit" &
        done

        wssh mkdir -p backup/lucasew/homelab/${config.networking.hostName}
        rsync -e "ssh -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i /run/secrets/rsyncnet-remote-backup" -avP /var/backup/ "${cfg.host}:backup/lucasew/homelab/${config.networking.hostName}"

        bash ${./hash-backups.sh}

        echo '[*] Waiting for jobs to finish...'

        while wait -n; do : ; done; # wait until it's possible to wait for bg job
      '';

      serviceConfig = {
        User = cfg.user;
        Group = cfg.group;
        ExecStartPre = [
        "+/run/current-system/sw/bin/chgrp -R ${cfg.group} /var/backup"
        "+/run/current-system/sw/bin/chmod -R g+r /var/backup"
        "+/run/current-system/sw/bin/find /var/backup -type d -exec /run/current-system/sw/bin/chmod g+x {} \\;"
        ];
      };
    };

    security.polkit.extraConfig = ''
      polkit.addRule(function(action, subject) {
          if (action.id == "org.freedesktop.systemd1.manage-units" && action.lookup("unit").startsWith('rsyncnet-remote-backup-git@') && subject.user === '${cfg.user}') {
            return polkit.Result.YES;
          }
      })
    '';

    systemd.services."rsyncnet-remote-backup-git@" = {
      path = [ wrappedSsh ];

      restartIfChanged = false;
      stopIfChanged = false;

      script = ''
        echo | wssh git --git-dir "git/$(echo "$INSTANCE_NAME" | sed 's;\/;-;g')" fetch --all --prune || true
      '';

      environment = {
        INSTANCE_NAME = "%I";
      };

      serviceConfig = {
        Type = "oneshot";
        TimeoutStartSec = cfg.git-step-timeout;
        User = cfg.user;
        Group = cfg.group;
      };
    };
  };
}
