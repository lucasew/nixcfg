{ config, pkgs, lib, ... }:

let
  cfg = config.programs.sunshine;
in

{
  options = {
    programs.sunshine = {
      enable = lib.mkEnableOption "sunshine";
      package = lib.mkPackageOption pkgs "sunshine" {};
    };
  };
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];

    systemd.user.services.sunshine = {
      script = ''
        xrandr --output $(xrandr  | grep [^s]connected | sed 's; ;\n;g'| head -n 1) --mode 1368x768
        ${lib.getExe cfg.package}
      '';
      serviceConfig = {
        Restart = "on-failure";
      };
    };
  };
}
