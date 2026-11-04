{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.programs.omarchy;
  omarchyIconFont = pkgs.runCommand "omarchy-icon-font" { } ''
    mkdir -p $out/share/fonts/truetype
    cp ${cfg.package}/share/omarchy/default/fonts/omarchy/omarchy.ttf \
      $out/share/fonts/truetype/omarchy.ttf
  '';
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
      pkgs.yaru-theme
    ];

    environment.sessionVariables.OMARCHY_PATH = "${cfg.package}/share/omarchy";

    # Bar glyphs are Nerd Font / omarchy.ttf PUA. Stock monospace here is
    # DejaVu, which has none of those codepoints.
    fonts.packages = [
      pkgs.nerd-fonts.jetbrains-mono
      pkgs.liberation_ttf
      omarchyIconFont
    ];
    fonts.fontconfig.defaultFonts.monospace = lib.mkBefore [ "JetBrainsMono Nerd Font" ];

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
        pkgs.fontconfig
      ]
      ++ lib.optional config.programs.sway.enable pkgs.sway
      ++ lib.optional config.programs.hyprland.enable config.programs.hyprland.package;
      environment = {
        OMARCHY_PATH = "${cfg.package}/share/omarchy";
      };
      script = ''
        export XDG_DATA_DIRS="${pkgs.yaru-theme}/share:${pkgs.adwaita-icon-theme}/share''${XDG_DATA_DIRS:+:$XDG_DATA_DIRS}"
        if [ -z "''${SWAYSOCK:-}" ]; then
          for sock in "''${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"/sway-ipc.*.sock; do
            if [ -S "$sock" ]; then
              export SWAYSOCK="$sock"
              break
            fi
          done
        fi
        exec ${lib.getExe' cfg.package "omarchy-launch-shell"}
      '';
    };
  };
}
