{
  config,
  lib,
  pkgs,
  ...
}:

{
  imports = [
    ../optional/flatpak.nix
    ../optional/kdeconnect-indicator.nix
    ../optional/dunst.nix
  ];

  config = lib.mkIf config.programs.hyprland.enable {

    services.xserver.displayManager.lightdm.enable = true;

    security.polkit.agent.enable = true;

    programs.hyprland = {
      xwayland.enable = true;
      # portalPackage = pkgs.xdg-desktop-portal-wlr // {
      #   override = args: pkgs.xdg-desktop-portal-wlr.override (builtins.removeAttrs args [ "hyprland" ]);
      # };
    };

    services.xserver.enable = true;

    services.dunst.enable = true;
    services.gammastep.enable = true;
    programs.waybar.enable = true;
    programs.kdeconnect.enable = true;

    # https://github.com/loki-47-6F-64/sunshine/commit/ebf9dbe9318808a5e127d3b6e397b9fa5149f197.patch
    # programs.sunshine.package = pkgs.sunshine.overrideAttrs (old: {
    #   patches = (old.patches or []) ++ [ ./sunshine-wayland.patch ];
    # });

    systemd.user.services.nm-applet = {
      path = with pkgs; [ networkmanagerapplet ];
      script = "exec nm-applet";
      restartIfChanged = true;
    };
    systemd.user.services.blueberry-tray = {
      path = with pkgs; [
        blueberry
        # (blueberry.overrideAttrs (old: {
        #   patches = (old.patches or [ ]) ++ [ ./blueberry-tray-fix.patch ];
        #   buildInputs = old.buildInputs ++ [ pkgs.libappindicator-gtk3 ];
        # }))
      ];
      script = "blueberry-tray; while true; do sleep 3600; done";
      restartIfChanged = true;
    };

    systemd.user.services.swayidle = {
      partOf = [ "graphical-session.target" ];
      path = with pkgs; [
        swayidle
        swaylock
        config.programs.hyprland.package
        playerctl
      ];
      restartIfChanged = true;
      script =
        with pkgs.custom.colors.colors;
        let
          swaylock-dict-args = {
            color = base00; # background

            text-color = base05;
            text-clear-color = base05;
            text-caps-lock-color = base05;
            text-ver-color = base05;
            text-wrong-color = base05;
            layout-text-color = base05;

            ring-color = base01;
            ring-clear-color = base0D;
            ring-caps-lock-color = base0C;
            ring-ver-color = base0A;
            ring-wrong-color = base08;

            key-hl-color = base06;
            bs-hl-color = base08;

            inside-color = "00000000";
            inside-clear-color = "00000000";
            inside-caps-lock-color = "00000000";
            inside-ver-color = "00000000";
            inside-wrong-color = "00000000";
            line-color = "00000000";
            line-clear-color = "00000000";
            line-caps-lock-color = "00000000";
            line-ver-color = "00000000";
            line-wrong-color = "00000000";
            layout-bg-color = "00000000";
            layout-border-color = "00000000";
          };
          swaylock-list-args = lib.pipe swaylock-dict-args [
            (builtins.mapAttrs (k: v: "--${k} ${v}"))
            (builtins.attrValues)
          ];
        in
        ''
          exec swayidle -w -d \
            idlehint 600 \
            timeout 605 'hyprctl dispatch dpms off' resume 'hyprctl dispatch dpms on' \
            lock 'swaylock -f ${lib.concatStringsSep " " swaylock-list-args}' \
            unlock 'hyprctl dispatch dpms on' \
            before-sleep 'playerctl pause'
        '';
    };

    security.pam.services.swaylock = { };
    environment.systemPackages = with pkgs; [
      swaylock
      eog # eye of gnome
      xfce.ristretto
      pcmanfm
      playerctl
      brightnessctl
      gscreenshot
      xwaylandvideobridge
      wl-clipboard
      custom.rofi_wayland
    ];
  };
}
