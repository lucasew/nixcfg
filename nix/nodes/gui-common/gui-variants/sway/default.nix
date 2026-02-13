{
  config,
  pkgs,
  lib,
  ...
}:

let
  lockerSpace = pkgs.makeDesktopItem {
    name = "locker";
    desktopName = "Bloquear Tela";
    icon = "lock";
    type = "Application";
    exec = "sdw utils i3wm lock-screen";
  };
  locker-params = let
    dict-args =
      with pkgs.custom.colors.colors; {

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

    swaylock-list-args = lib.pipe dict-args [
      (builtins.mapAttrs (k: v: ["--${k}" "${v}"]))
      (builtins.attrValues)
      (lib.flatten)
    ];
  in swaylock-list-args;
in

{
  imports = [
    ../optional/flatpak.nix
    ../optional/kdeconnect-indicator.nix
    ../optional/dunst.nix
    ../../workspaced.nix
  ];
  config = lib.mkIf config.programs.sway.enable {
    systemd.user.services.xss-lock.restartIfChanged = true;

    programs.ssh.startAgent = true;

    security.soteria.enable = true;

    programs.xss-lock = {
      enable = true;
      lockerCommand = lib.mkDefault ''
        ${lib.getExe pkgs.swaylock} ${lib.escapeShellArgs locker-params}
      '';
    };

    # Swayidle for power management
    systemd.user.services.swayidle = {
      partOf = [ "graphical-session.target" ];
      path = with pkgs; [ swayidle procps ];
      restartIfChanged = true;
      script = ''
        PATH=$PATH:/run/current-system/sw/bin
        exec swayidle -w -d \
          timeout 300 'workspaced system power lock' \
          timeout 10 'pgrep swaylock && workspaced system screen off' \
          resume 'workspaced system screen on' \
          before-sleep 'workspaced system power lock'
      '';
    };

    # XDG Portal for Wayland
    xdg.portal = {
      enable = true;
      xdgOpenUsePortal = true;
      extraPortals = [ pkgs.xdg-desktop-portal-wlr ];
      config.common.default = "wlr";
      wlr.enable = true;
    };

    # Display Manager
    services.displayManager.sessionData.autologinSession = lib.mkDefault "sway";
    services.xserver.displayManager.lightdm.enable = true;
    services.xserver.enable = true;

    # System tray services
    systemd.user.services.nm-applet = {
      path = with pkgs; [ networkmanagerapplet ];
      script = "nm-applet";
    };
    systemd.user.services.blueberry-tray = {
      path = with pkgs; [ blueberry ];
      script = "blueberry-tray; while true; do sleep 3600; done";
    };

    services.flatpak.enable = true;
    services.tumbler.enable = true;
    services.dunst.enable = true;
    services.gammastep.enable = true;
    programs.waybar.enable = true;
    programs.kdeconnect.enable = true;

    # System packages
    environment.systemPackages = with pkgs; [
      wlr-randr
      wl-clipboard
      grim
      imv
      eog # eye of gnome
      xfce.ristretto
      pcmanfm
      kitty
      slurp
      rofi
      swaybg
      lockerSpace
      playerctl
      pulseaudio
      feh
      brightnessctl
      unstable.i3pystatus
    ];

  };
}
