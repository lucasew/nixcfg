{
  config,
  pkgs,
  lib,
  ...
}:

let
  cfg = config.services.sunshine;
in

{
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];

    services.sunshine.settings = {
      gamepad = "ds4";
    };

    systemd.user.services.sunshine = {
      serviceConfig = {
        Restart = "on-failure";
      };
    };
  };
}
