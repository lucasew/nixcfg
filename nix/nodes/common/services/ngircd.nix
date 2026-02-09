{
  config,
  pkgs,
  lib,
  ...
}: let
  cfg = config.services.ngircd;

  toml = pkgs.formats.toml {};

  configFile =
    pkgs.runCommand "ngircd.conf"
    {
      config = toml.generate "ngircd_input.conf" cfg.config;

      preferLocalBuild = true;
    }
    ''
      cp $config $out

      # TODO: proper generator
      substituteInPlace $out \
        --replace '= false' '= "no"' \
        --replace '= true' '= "yes"'
      sed -i 's;^\[\[\([^\]*\)\]\]$;[\1];' $out # general, for example, may appear twice, fixing syntax
      sed -i 's;^\([^ ]*\) = \"\([^"]*\)"$;\1 = \2;' $out # fix quote enclosing
      sed -i 's;^\([^=]*\)=;  \1=;' $out


      ${lib.getExe cfg.package} --config $out --configtest
    '';
in {
  disabledModules = ["services/networking/ngircd.nix"];

  options = {
    services.ngircd = {
      enable = lib.mkEnableOption "the ngircd IRC server";

      config = lib.mkOption {
        description = "The ngircd configuration (see ngircd.conf(5)).";

        type = toml.type;
      };

      package = lib.mkPackageOption pkgs "ngircd" {};
    };
  };

  config = lib.mkIf cfg.enable {
    #!!! TODO: Use ExecReload (see https://github.com/NixOS/nixpkgs/issues/1988)
    environment.etc."ngircd.conf".source = configFile;

    systemd.services.ngircd = {
      description = "The ngircd IRC server";

      wantedBy = ["multi-user.target"];

      serviceConfig = {
        ExecStart = "${lib.getExe cfg.package} --config /etc/ngircd.conf --nodaemon --syslog";
        ExecReload = "/run/current-system/sw/bin/kill -HUP $MAINPID";
        User = "ngircd";
        Group = "ngircd";
        PrivateTmp = true;
        NoNewPrivileges = true;
        PrivateDevices = true;
        DevicePolicy = "closed";
        ProtectSystem = "strict";
        ProtectHome = true;
        ProtectControlGroups = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
        MemoryDenyWriteExecute = true;
        LockPersonality = true;
      };

      # systemd sends SIGHUP on reload, which is supported
      reloadIfChanged = true;
    };

    users.users.ngircd = {
      isSystemUser = true;
      group = "ngircd";
      description = "ngircd user.";
    };
    users.groups.ngircd = {};

    # Logic related to my config
    # TODO: did I read socket activation????
    networking.ports.ircd.enable = true;
    services.ngircd.package = lib.mkDefault pkgs.unstable.ngircd;
    services.ngircd.config = {
      Global = {
        Info = "lucasew's IRC server";
        Listen = "127.0.0.1"; # ts-proxy will reverse proxy it
        MotdPhrase = "aoba";
        Ports = config.networking.ports.ircd.port;
      };
      Channel = [
        {
          Name = "#general";
          Topic = "Tópico principal";
          Autojoin = true;
        }
        {
          Name = "#test";
          Topic = "Tópico teste";
        }
      ];
      Options = {
        PAM = false;
      };
      Limits = {
        MaxNickLength = 16;
      };
    };
    services.ts-proxy.hosts.irc = {
      enableTLS = true;
      enableRaw = true;
      address = "127.0.0.1:${toString config.networking.ports.ircd.port}";
      proxies = ["ngircd.service"];
      listen = 6697;
    };
  };
}
