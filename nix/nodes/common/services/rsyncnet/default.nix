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
        extraGroups = [ "ssh" ];
      };
      groups.${cfg.group} = { };
    };

    sops.secrets.rsyncnet-remote-backup = {
      sopsFile = ../../../../../secrets/rsyncnet;
      owner = cfg.user;
      group = cfg.group;
      format = "binary";
    };

    systemd.tmpfiles.rules = [ "d ${cfg.dataDir} 0700 ${cfg.user} ${cfg.group} - -" ];

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
      path = with pkgs; [
        wrappedSsh
        bash
        pv
        git
        gawk
        rsync
        openssh
      ];

      restartIfChanged = false;
      stopIfChanged = false;

      script = ''
        cd "${cfg.dataDir}"
        export PATH+=":/run/wrappers/bin"

        wssh mkdir -p backup/lucasew/homelab/${config.networking.hostName}
        # cópia dos backups do postgres por ex
        rsync -e "ssh -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i /run/secrets/rsyncnet-remote-backup" -avP /var/backup/ "${cfg.host}:backup/lucasew/homelab/${config.networking.hostName}"

        function backup_git {
          repo="$1"; shift
          echo git-backup $repo >&2
          timeout 60 wssh git --git-dir "git/$repo" fetch --all --prune || {
            printf "Subject: git-backup/falha: %s\n%s" "$repo" "Backup do repositório falhou" | sendmail
          }
        }

        for repo in $(wssh ls git | grep -v -e '^zzz'); do
          backup_git "$repo"
        done

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
  };
}
