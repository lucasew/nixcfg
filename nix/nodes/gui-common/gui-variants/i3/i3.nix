{
  pkgs,
  lib,
  config,
  ...
}:

let
  custom_rofi = pkgs.custom.rofi.override { inherit (pkgs.custom) colors; };
  inherit (pkgs) makeDesktopItem;
  mod = "Mod4";
  lockerSpace = makeDesktopItem {
    name = "locker";
    desktopName = "Bloquear Tela";
    icon = "lock";
    type = "Application";
    exec = "sdw utils i3wm lock-screen";
  };
in
{

  config = lib.mkIf config.services.xserver.windowManager.i3.enable {
    # https://thomashunter.name/i3-configurator/
    environment.systemPackages = [
      lockerSpace
      custom_rofi
      pkgs.kitty
    ];
    services.xserver.displayManager.lightdm = {
      greeters.enso = {
        enable = true;
        blur = true;
      };
    };

    services.xserver.windowManager.i3 = {
      configFile = "/etc/i3config";
      extraPackages = with pkgs; [
        playerctl
        custom_rofi
        pulseaudio
        feh
        brightnessctl
      ];
    };
    environment.etc."i3status".text =
      with pkgs.custom.colors.colors;
      lib.mkForce ''
        general {
          colors = true
          color_bad = "#${base08}"
          color_degraded = "#${base08}"
          color_good = "#${base0B}"

          interval = 2
        }

        order += "read_file netusage"
        read_file netusage {
          path = "/dev/shm/netusage"
          format = "💬 %content"
          format_bad = ""
        }

        order += "wireless _first_"
        wireless _first_ {
          format_up = "🛜 %bitrate %frequency %quality"
          format_down = "🛜 OFF"
        }

        order += "ethernet _first_"
        ethernet _first_ {
                format_up = "🔗 %speed"
                format_down = "🔗 OFF"
        }

        order += "load"
        load {
                format = "💪 %1min %5min %15min"
        }

        order += "cpu_usage"
        cpu_usage {
                format= "🧠 %usage"
                max_threshold= 75
        }

        order += "memory"
        memory {
              format = "🐸 %percentage_used"
              threshold_degraded = "10%"
              format_degraded = "🐸 %free"
        }

        order += "volume master"
        volume master {
                format = "🔊 %volume"
                format_muted = "🔇"
                device = "default"
                mixer = "Master"
                mixer_idx = 0
        }

        order += "battery all"
        battery all {
                format = "%status %percentage"
                format_down = ""
                status_chr = "⚡"
                status_bat = "🔋"
                status_unk = "🤔"
                status_full = "🔌"
                path = "/sys/class/power_supply/BAT%d/uevent"
                low_threshold = 10
        }

        order += "tztime local"
        tztime local {
                format = "⌛ %Y-%m-%d %H:%M"
        }


      '';
    environment.etc."i3config".text =
      with pkgs.custom.colors.colors;
      lib.mkForce ''
        set $mod ${mod}

        bar {
          status_command ${lib.getExe pkgs.i3pystatus} -c $(sdw d root)/bin/_shortcuts/i3pystatus/main.py
          i3bar_command i3bar --transparency
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
        bindsym $mod+Down focus down
        bindsym $mod+Left focus left
        bindsym $mod+Return exec kitty
        bindsym $mod+Right focus right
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

        bindsym $mod+Shift+minus move scratchpad
        bindsym $mod+Shift+q kill
        bindsym $mod+Shift+r restart
        bindsym $mod+Up focus up
        bindsym $mod+a focus parent
        bindsym $mod+d exec rofi-launch
        bindsym $mod+Shift+d exec rofi-window
        bindsym $mod+e layout toggle split

        bindsym $mod+f fullscreen toggle
        bindsym $mod+Ctrl+f fullscreen toggle global

        bindsym $mod+minus scratchpad show
        bindsym $mod+s layout stacking
        bindsym $mod+space focus mode_toggle

        bindsym $mod+h split h
        bindsym $mod+v split v

        bindsym $mod+w layout tabbed
        bindsym $mod+Ctrl+Right=resize shrink width 1 px or 1 ppt
        bindsym $mod+Ctrl+Up=resize grow height 1 px or 1 ppt
        bindsym $mod+Ctrl+Down=resize shrink height 1 px or 1 ppt
        bindsym $mod+Ctrl+Left=resize grow width 1 px or 1 ppt

        # custom keys
        bindsym XF86AudioRaiseVolume exec pactl set-sink-volume @DEFAULT_SINK@ +10%
        bindsym XF86AudioLowerVolume exec pactl set-sink-volume @DEFAULT_SINK@ -10%
        bindsym XF86AudioMute exec pactl set-sink-volume @DEFAULT_SINK@ toggle
        bindsym XF86AudioMicMute exec pactl set-sink-volume @DEFAULT_SOURCE@ toggle
        bindsym XF86AudioNext exec playerctl next
        bindsym XF86AudioPrev exec playerctl previous
        bindsym XF86AudioPlay exec playerctl play-pause
        bindsym XF86AudioPause exec playerctl play-pause
        bindsym XF86MonBrightnessUp exec brightnessctl s +5%
        bindsym XF86MonBrightnessDown exec brightnessctl s 5%-

        bindsym $mod+Shift+m move workspace to output left
        bindsym $mod+m focus output next

        bindsym $mod+l exec sdw utils i3wm lock-screen
        bindsym $mod+n exec sdw utils i3wm modn

        bindsym $mod+b exec sdw utils i3wm goto-new-ws
        bindsym $mod+Shift+b exec sdw utils i3wm goto-new-ws shift


        bindsym --release Print exec org.flameshot.Flameshot gui

        exec --no-startup-id ${pkgs.mate.mate-polkit}/libexec/polkit-mate-authentication-agent-1

        # exec_always feh --bg-fill --no-xinerama --no-fehbg '/etc/wallpaper'
        exec_always feh --bg-fill --no-fehbg '/etc/wallpaper'

        exec_always systemctl restart --user nm-applet.service blueberry-tray.service kdeconnect.service kdeconnect-indicator.service

        default_border pixel 2
        hide_edge_borders smart
        focus_on_window_activation urgent
      '';
  };
}
