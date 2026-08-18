{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.programs.omarchy;
in
{
  options.programs.omarchy = {
    enable = lib.mkEnableOption "Omarchy CLI, plugin manager, and optional Quickshell host";

    package = lib.mkPackageOption pkgs "omarchy" { };

    shell.enable = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = ''
        Start omarchy-shell (Quickshell) with the graphical session.
        That process is what `omarchy plugin add` loads into.
        On Sway you will have two bars until Waybar is turned off.
        Workspace widgets that import Quickshell.Hyprland do not work on Sway.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [
      cfg.package
      pkgs.quickshell
    ];

    environment.sessionVariables.OMARCHY_PATH = "${cfg.package}/share/omarchy";

    systemd.user.services.omarchy-shell = lib.mkIf cfg.shell.enable {
      description = "Omarchy Quickshell desktop host";
      partOf = [ "graphical-session.target" ];
      after = [ "graphical-session.target" ];
      wantedBy = [ "graphical-session.target" ];
      restartIfChanged = true;
      path = [
        cfg.package
        pkgs.quickshell
        pkgs.jq
        pkgs.gum
        pkgs.git
      ]
      ++ lib.optional config.programs.sway.enable pkgs.sway
      ++ lib.optional config.programs.hyprland.enable config.programs.hyprland.package;
      environment = {
        OMARCHY_PATH = "${cfg.package}/share/omarchy";
      };
      script = "exec ${lib.getExe' cfg.package "omarchy-launch-shell"}";
    };
  };
}
