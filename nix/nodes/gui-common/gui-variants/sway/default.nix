{
  config,
  pkgs,
  lib,
  ...
}:

let
  custom_rofi = pkgs.custom.rofi_wayland.override { inherit (pkgs.custom) colors; };
  mod = "Mod4";
  lockerSpace = pkgs.makeDesktopItem {
    name = "locker";
    desktopName = "Bloquear Tela";
    icon = "lock";
    type = "Application";
    exec = "sdw utils i3wm lock-screen";
  };
in

{
  imports = [
    ../optional/flatpak.nix
    ../optional/kdeconnect-indicator.nix
    ../optional/dunst.nix
  ];
  config = lib.mkIf config.programs.sway.enable {
    services.dunst.enable = true;
    xdg.portal = {
      enable = true;
      config = {
        common = {
          default = "wlr";
        };
      };
      wlr.enable = true;
      # wlr.settings.screencast = { output_name = "DP-2"; chooser_type = "simple"; chooser_cmd = "${pkgs.slurp}/bin/slurp -f %o -or"; };
    };
    services.xserver.displayManager.lightdm.enable = true;
    services.xserver.displayManager.sessionData.autologinSession = lib.mkDefault "sway";
    services.xserver.enable = true;
    systemd.user.services.gammastep.environment.WAYLAND_DISPLAY = "wayland-1";

    security.polkit.agent.enable = true;
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
      eog # eye of gnome
      xfce.ristretto
      pcmanfm
      kitty
      custom_rofi
      lockerSpace
      playerctl
      pulseaudio
      feh
      brightnessctl
    ];

    systemd.user.services.xss-lock.restartIfChanged = true;

    environment.etc."sway/config".text =
      with pkgs.custom.colors.colors;
      lib.mkForce ''
        set $mod ${mod}

        input type:keyboard {
          xkb_layout br,us
          xkb_options grp:win_space_toggle,terminate:ctrl_alt_bksp
        }

        bar {
          status_command ${lib.getExe pkgs.unstable.i3pystatus} -c $(sdw d root)/bin/_shortcuts/i3pystatus/main.py
          # i3bar_command i3bar --transparency
          font pango: Fira Code 10
          hidden_state show
          position top
          # output primary
          tray_output primary
          workspace_buttons yes

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
        bindsym $mod+Return exec sdw shim terminal

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

        bindsym $mod+minus exec sdw utils i3wm toggle-scratchpad
        bindsym $mod+Shift+minus move scratchpad

        bindsym $mod+Shift+q kill
        bindsym $mod+Shift+r restart
        bindsym $mod+d exec rofi-launch
        bindsym $mod+Shift+d exec rofi-window
        bindsym $mod+e layout toggle split

        bindsym $mod+f fullscreen toggle
        bindsym $mod+Ctrl+f fullscreen toggle global

        bindsym $mod+s layout stacking
        bindsym $mod+space focus mode_toggle

        bindsym $mod+h split h
        bindsym $mod+v split v

        bindsym $mod+w layout tabbed
        bindsym $mod+Ctrl+Right resize shrink width 1 px or 1 ppt
        bindsym $mod+Ctrl+Up resize grow height 1 px or 1 ppt
        bindsym $mod+Ctrl+Down resize shrink height 1 px or 1 ppt
        bindsym $mod+Ctrl+Left resize grow width 1 px or 1 ppt

        # custom keys
        bindsym XF86AudioRaiseVolume exec sdw utils i3wm audio up
        bindsym XF86AudioLowerVolume exec sdw utils i3wm audio down
        bindsym XF86AudioMute exec sdw utils i3wm audio mute

        bindsym XF86AudioNext exec  sdw utils i3wm playerctl next
        bindsym XF86AudioPrev exec  sdw utils i3wm playerctl previous
        bindsym XF86AudioPlay exec  sdw utils i3wm playerctl play-pause
        bindsym XF86AudioPause exec sdw utils i3wm playerctl play-pause

        bindsym XF86MonBrightnessUp   exec sdw utils i3wm brightnessctl up
        bindsym XF86MonBrightnessDown exec sdw utils i3wm brightnessctl down

        bindsym $mod+Shift+m move workspace to output left
        bindsym $mod+m focus output next

        bindsym $mod+l exec sdw utils i3wm lock-screen
        bindsym $mod+n exec sdw utils i3wm modn

        bindsym $mod+b exec sdw utils i3wm goto-new-ws
        bindsym $mod+Shift+b exec sdw utils i3wm goto-new-ws window


        bindsym --release Print exec org.flameshot.Flameshot gui

        exec --no-startup-id ${pkgs.mate.mate-polkit}/libexec/polkit-mate-authentication-agent-1

        # exec_always feh --bg-fill --no-xinerama --no-fehbg '/etc/wallpaper'
        # exec_always feh --bg-fill --no-fehbg '/etc/wallpaper'

        exec_always systemctl restart --user nm-applet.service blueberry-tray.service kdeconnect.service kdeconnect-indicator.service

        default_border pixel 2
        hide_edge_borders smart
        focus_on_window_activation urgent
      '';
  };
}
