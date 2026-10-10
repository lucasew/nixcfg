{
  config,
  pkgs,
  lib,
  self,
  ...
}:
let
  cfg = config.programs.kdeconnect;
  script = builtins.readFile "${self}/bin/waybar/kdeconnect";
  waybarKdeconnect = pkgs.writeShellApplication {
    name = "waybar-kdeconnect";
    runtimeInputs = [
      pkgs.systemd
      pkgs.coreutils
      pkgs.gnugrep
      cfg.package
    ];
    text = lib.removePrefix "#!/usr/bin/env bash\n" script;
  };
in
{
  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ waybarKdeconnect ];

    systemd.user.services.kdeconnect = {
      enable = true;
      script = "${cfg.package}/libexec/kdeconnectd";
      restartIfChanged = true;
    };
    systemd.user.services.kdeconnect-indicator = {
      enable = true;
      path = [ cfg.package ];
      script = "kdeconnect-indicator";
      restartIfChanged = true;
    };

    # waybar.service PATH is a store subset; without this, exec/on-click miss kdeconnect-cli.
    systemd.user.services.waybar = lib.mkIf config.programs.waybar.enable {
      path = [
        waybarKdeconnect
        cfg.package
      ];
    };
  };
}
