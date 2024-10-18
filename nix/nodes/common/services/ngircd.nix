{ config, pkgs, lib, ...}:


let
  cfg = config.services.ngircd;

  toml = pkgs.formats.toml {};

  configFile = pkgs.runCommand "ngircd.conf" {
    config = toml.generate "ngircd_input.conf" cfg.config;

    preferLocalBuild = true;
  } ''
      cp $config $out
      ${cfg.package}/sbin/ngircd --config $out --configtest
  '';

in
{
  disabledModules = [ "services/networking/ngircd.nix" ];  

  options = {
    services.ngircd = {
      enable = lib.mkEnableOption "the ngircd IRC server";

      config = lib.mkOption {
        description = "The ngircd configuration (see ngircd.conf(5)).";

        type = toml.type;
      };

      package = lib.mkPackageOption pkgs "ngircd" { };
    };
  };

  config = lib.mkIf cfg.enable {
    #!!! TODO: Use ExecReload (see https://github.com/NixOS/nixpkgs/issues/1988)
    systemd.services.ngircd = {
      description = "The ngircd IRC server";

      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        ExecStart = "${lib.getExe cfg.package} --config ${configFile} --nodaemon";
        User = "ngircd";
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
    services.ngircd.config = {
      global = {
        Info = "lucasew's IRC server";
        Listen = "127.0.0.1"; # ts-proxy will reverse proxy it
        MotdPhrase = "aoba";
        Ports = config.networking.ports.ircd.port;
      };
      channel = [
        {
          name = "#general";
          topic = "Tópico principal";
        }
        {
          name = "#test";
          topic = "Tópico teste";
        }
      ];
    };
    services.ts-proxy.hosts.irc = {
      enableTLS = true;
      enableRaw = true;
      address = "127.0.0.1:${toString config.networking.ports.ircd.port}";
      listen = 6697;
    };
  };
}
