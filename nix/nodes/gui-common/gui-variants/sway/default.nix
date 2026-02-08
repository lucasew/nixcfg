{
  config,
  pkgs,
  lib,
  ...
}:

let
  mod = "Mod4";
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

    systemd.user.services.swayidle = {
      partOf = [ "graphical-session.target" ];
      path = with pkgs; [ swayidle procps ];
      restartIfChanged = true;
      script = ''
        PATH=$PATH:/run/current-system/sw/bin
        exec swayidle -w -d \
          timeout 300 'loginctl lock-session' \
          timeout 10 'pgrep swaylock && swaymsg "output * dpms off"' \
          resume 'swaymsg "output * dpms on"' \
          before-sleep 'loginctl lock-session'
      '';
    };

    services.dunst.enable = true;
    xdg.portal = {
      enable = true;
      xdgOpenUsePortal = true;
      extraPortals = [
        pkgs.xdg-desktop-portal-wlr
      ];
      config = {
        common = {
          default = "wlr";
        };
      };
      wlr.enable = true;
      # wlr.settings.screencast = { output_name = "DP-2"; chooser_type = "simple"; chooser_cmd = "${pkgs.slurp}/bin/slurp -f %o -or"; };
    };
    services.displayManager.sessionData.autologinSession = lib.mkDefault "sway";
    services.xserver.displayManager.lightdm.enable = true;
    services.xserver.enable = true;
    # systemd.user.services.gammastep.environment.WAYLAND_DISPLAY = "wayland-1";

    services.tumbler.enable = true;
    programs.kdeconnect.enable = true;
    services.gammastep.enable = true;
    systemd.user.services.nm-applet = {
      path = with pkgs; [ networkmanagerapplet ];
      script = "nm-applet";
    };
    systemd.user.services.blueberry-tray = {
      path = with pkgs; [ blueberry ];
      script = "blueberry-tray; while true; do sleep 3600; done";
    };
    environment.systemPackages = with pkgs; [
      wlr-randr
      wl-clipboard
      eog # eye of gnome
      xfce.ristretto
      pcmanfm
      kitty
      slurp
      rofi-wayland
      lockerSpace
      playerctl
      pulseaudio
      feh
      brightnessctl
    ];

    programs.waybar.enable = true;

    environment.etc."sway/config".text =
      with pkgs.custom.colors.colors;
      lib.mkForce ''
        set $mod ${mod}

        input type:keyboard {
          xkb_layout br,us
          xkb_options grp:win_space_toggle,terminate:ctrl_alt_bksp
        }

        bar {
          # swaybar_command ${lib.getExe pkgs.waybar}
          status_command ${lib.getExe pkgs.unstable.i3pystatus} -c $(sdw d root)/bin/_shortcuts/i3pystatus/main.py
          font pango: Fira Code 10
          hidden_state hide
          mode invisible
          position top
          # output primary
          # tray_output primary
          # workspace_buttons yes

          colors {
            # background #00${base00}
            background #00000000
            statusline #${base05}
            separator #${base00}

            # name             border     background text
            focused_workspace  #${base01} #${base02} #${base05}
            active_workspace   #${base01} #${base03} #${base05}
            inactive_workspace #${base01} #${base01} #${base05}
            urgent_workspace   #${base08} #${base08} #${base00}
            binding_mode       #${base00} #${base00} #${base05}

          }
        }

        # Property Name         Border    Background Text     Indicator  Child
        client.focused          #${base01} #${base00} #${base05} #${base0D} #${base0C}
        client.focused_inactive #${base01} #${base01} #${base05} #${base03} #${base01}
        client.unfocused        #${base01} #${base02} #${base05} #${base01} #${base01}
        client.urgent           #${base08} #${base08} #${base00} #${base08} #${base08}
        client.placeholder      #${base00} #${base00} #${base05} #${base00} #${base00}
        client.background       #${base07} #${base00} #${base05}

        bindsym $mod+0 workspace number 10
        bindsym $mod+1 workspace number 1
        bindsym $mod+2 workspace number 2
        bindsym $mod+3 workspace number 3
        bindsym $mod+4 workspace number 4
        bindsym $mod+5 workspace number 5
        bindsym $mod+6 workspace number 6
        bindsym $mod+7 workspace number 7
        bindsym $mod+8 workspace number 8
        bindsym $mod+9 workspace number 9

        bindsym $mod+Shift+0 move container to workspace number 10
        bindsym $mod+Shift+1 move container to workspace number 1
        bindsym $mod+Shift+2 move container to workspace number 2
        bindsym $mod+Shift+3 move container to workspace number 3
        bindsym $mod+Shift+4 move container to workspace number 4
        bindsym $mod+Shift+5 move container to workspace number 5
        bindsym $mod+Shift+6 move container to workspace number 6
        bindsym $mod+Shift+7 move container to workspace number 7
        bindsym $mod+Shift+8 move container to workspace number 8
        bindsym $mod+Shift+9 move container to workspace number 9

        bindsym $mod+Down focus down
        bindsym $mod+Up focus up
        bindsym $mod+Left focus left
        bindsym $mod+Right focus right
        bindsym $mod+a focus parent
        bindsym $mod+Return exec workspaced launch terminal

        bindsym $mod+Shift+Down move down
        bindsym $mod+Shift+Left move left
        bindsym $mod+Shift+Right move right
        bindsym $mod+Shift+Up move up

        bindsym $mod+Shift+h workspace prev_on_output
        bindsym $mod+Shift+l workspace next_on_output

        bindsym $mod+Shift+c reload
        bindsym $mod+Shift+e exec i3-nagbar -t warning -m 'Do you want to exit i3?' -b 'Yes' 'loginctl kill-session $XDG_SESSION_ID'
        bindsym $mod+Shift+f floating toggle
        bindsym $mod+Shift+s sticky toggle

        bindsym $mod+minus exec workspaced dispatch workspace scratchpad
        bindsym --locked $mod+equal exec workspaced dispatch screen toggle
        bindsym $mod+Shift+minus move scratchpad

        bindsym $mod+Shift+q kill
        bindsym $mod+Shift+r restart
        bindsym $mod+d exec workspaced dispatch menu launch
        bindsym $mod+Shift+d exec workspaced dispatch menu window
        bindsym $mod+e layout toggle split

        bindsym $mod+f fullscreen toggle
        bindsym $mod+Ctrl+f fullscreen toggle global

        bindsym $mod+s layout stacking
        bindsym $mod+space focus mode_toggle

        bindsym $mod+h split h
        bindsym $mod+v split v

        bindsym $mod+w layout tabbed
        bindsym $mod+Ctrl+Right resize shrink width 1 px or 1 ppt
        bindsym $mod+Ctrl+Up    resize grow height 1 px or 1 ppt
        bindsym $mod+Ctrl+Down  resize shrink height 1 px or 1 ppt
        bindsym $mod+Ctrl+Left  resize grow width 1 px or 1 ppt

        # custom keys
        bindsym --locked XF86AudioRaiseVolume exec workspaced dispatch audio up
        bindsym --locked XF86AudioLowerVolume exec workspaced dispatch audio down
        bindsym --locked XF86AudioMute exec workspaced dispatch audio mute

        bindsym --locked XF86AudioNext exec  workspaced dispatch media next
        bindsym --locked XF86AudioPrev exec  workspaced dispatch media previous
        bindsym --locked XF86AudioPlay exec  workspaced dispatch media play-pause
        bindsym --locked XF86AudioPause exec workspaced dispatch media play-pause

        bindsym --locked XF86MonBrightnessUp   exec workspaced dispatch brightness up
        bindsym --locked XF86MonBrightnessDown exec workspaced dispatch brightness down

        bindsym $mod+Shift+m move workspace to output left
        bindsym $mod+m focus output right

        bindsym --locked $mod+l exec workspaced dispatch screen lock
        bindsym $mod+n exec workspaced dispatch workspace rotate
        bindsym $mod+z exec workspaced dispatch menu workspace
        bindsym $mod+Shift+z exec workspaced dispatch menu workspace --move


        bindsym $mod+b exec workspaced dispatch workspace next
        bindsym $mod+Shift+b exec workspaced dispatch workspace next --move


        bindsym --release Print exec workspaced dispatch screenshot output
        bindsym $mod+Print exec workspaced dispatch screenshot select
        bindsym $mod+Shift+Print exec workspaced dispatch screenshot window
        bindsym $mod+Ctrl+Shift+Print exec workspaced dispatch screenshot all


        output * bg /etc/wallpaper fill

        # exec_always feh --bg-fill --no-xinerama --no-fehbg '/etc/wallpaper'
        exec_always feh --bg-fill --no-fehbg '/etc/wallpaper'

        default_border pixel 2
        hide_edge_borders smart
        focus_on_window_activation urgent

        

        # announce a running sway session to systemd
        # https://discourse.nixos.org/t/xdg-desktop-portal-not-working-on-wayland-while-kde-is-installed/20919/3
        exec_always systemctl --user import-environment XDG_SESSION_TYPE XDG_CURRENT_DESKTOP WAYLAND_DISPLAY DISPLAY SWAYSOCK PATH
        exec_always dbus-update-activation-environment --systemd WAYLAND_DISPLAY XDG_CURRENT_DESKTOP=sway

        exec_always systemctl restart --user xss-lock nm-applet.service blueberry-tray.service kdeconnect.service kdeconnect-indicator.service dunst.service swayidle.service waybar.service
        exec_always systemctl --user restart xdg-desktop-portal
        exec_always env WORKSPACED_REFRESH=1 workspaced
        exec_always workspaced dispatch screen reset
      '';
  };
}
