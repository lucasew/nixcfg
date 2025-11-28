{
  global,
  pkgs,
  lib,
  self,
  ...
}:
let
  inherit (lib.hm.gvariant) mkTuple;
  inherit (pkgs.custom) colors;
in
{

  imports = [
    ../base/default.nix
    ./atuin.nix
    ./dlna.nix
    ./helix
    ./ghostty.nix
    ./espanso.nix
    ./dconf.nix
    ./borderless-browser.nix
    ./theme
    ./discord.nix
    ./qutebrowser.nix
    ./zen-browser.nix
    ./mise.nix
  ];

  stylix.enable = true;

  borderless-browser.chromium = lib.getExe pkgs.vivaldi;

  # programs.ghostty.enable = true;

  programs.atuin.enable = true;

  programs.zen-browser.enable = true;
  programs.vscode.enable = true;
  programs.helix.enable = true;
  # services.espanso.enable = true;
  programs.man.enable = true;

  # programs.qutebrowser.enable = true;

  home = {
    homeDirectory = /home/lucasew;
    inherit (global) username;
  };

  home.packages = with pkgs; [
    unstable.zed-editor
    uv
    ruff
    mission-center
    # custom.firefox # now I am using chromium
    cached-nix-shell
    devenv
    dotenv
    jless # json viewer
    feh
    fortune
    graphviz
    github-cli
    google-cloud-sdk
    htop
    libnotify
    ncdu
    # nix-option
    nix-prefetch-scripts
    nix-output-monitor
    pkg
    rclone
    ripgrep
    fd
    remmina
    sqlite
    sshpass

    # media
    nbr.wine-apps._7zip
    xxd

    # custom.vscode.programming
    # (custom.neovim.override { inherit colors; })
    # (custom.emacs.override { inherit colors; })

    # LSPs
    nil
    python3Packages.python-lsp-server
    (pkgs.writeShellScriptBin "e" ''
      if [ ! -v EDITOR ]; then
        export EDITOR=hx
      fi
      "$EDITOR" "$@"
    '')
    (pkgs.makeDesktopItem {
      name = "nixcfg-quicksync";
      desktopName = "nixcfg: Sincronização Rápida";
      icon = "sync-synchronizing";
      exec = "sdw quicksync";
    })
    (pkgs.makeDesktopItem {
      name = "nixcfg-backup";
      desktopName = "nixcfg: Backup";
      icon = "sync-synchronizing";
      exec = "sdw backup";
    })
  ];

  # programs.hello-world.enable = true;

  programs = {
    # adskipped-spotify.enable = true;
    jq.enable = true;
    obs-studio = {
      package = pkgs.obs-studio;
      enable = true;
    };
    
  };

  gtk = {
    enable = true;
  };
  qt = {
    enable = true;
    platformTheme.name = "gtk";
  };

  programs.terminator = {
    # enable = true;
    config = {
      global_config.borderless = true;
    };
  };
  programs.bash.enable = true;

  programs.mpv = {
    enable = true;
    config = {
      ytdl-raw-options = "format-sort=\"vcodec:h264,res,acodec:m4a\"";
    };
  };

  programs.waybar = {
    enable = true;
    settings = {   
      settings = {
        layer = "top"; # Waybar at top layer
        position = "top"; # Waybar position (top|bottom|left|right)
        height = 10; # Waybar height (to be removed for auto height)
        width = 1280; # Waybar width
        spacing = 1; # Gaps between modules (4px)
        # modules-left = ["hyprland/workspaces"];
        modules-left = ["sway/workspaces" "sway/mode"];
        modules-center = ["hyprland/window"];
        modules-right = [
          "idle_inhibitor"
          "pulseaudio"
          "backlight"
          "network"
          "custom/updates"
          "cpu"
          "memory"
          "temperature"
          "battery"
          "tray"
          "clock"
        ];
        keyboard-state = {
            numlock = false;
            capslock = false;
            format = "{name} {icon}";
            format-icons = {
                locked = "🔒";
                unlocked = "🔓";
            };
        };
        "hyprland/window" = {
            max-length = 50;
            separate-outputs = true;
        };
        idle_inhibitor = {
            format = "{icon}";
            format-icons = {
                activated = "☕";
                deactivated = "💤";
            };
        };
        tray = {
            # "icon-size": 21,
            spacing = 5;
        };
        clock = {
            timezone = "America/Sao_Paulo";
            tooltip-format = "<big>{:%Y %B}</big>\n<tt><small>{calendar}</small></tt>";
            format-alt = "{:%Y-%m-%d}";
        };
        cpu = {
            format = "{usage}%";
            tooltip = false;
        };
        memory = {
            format = "{}%";
        };
        temperature = {
            # "thermal-zone": 2,
            # "hwmon-path": "/sys/class/hwmon/hwmon2/temp1_input",
            critical-threshold = 70;
            format-critical = "{temperatureC}°C!";
            format = "{temperatureC}°C";
        };
        battery = {
            states = {
                # "good": 95,
                warning = 30;
                critical = 15;
            };
            format = "{capacity}% {icon}";
            format-charging = "{capacity}% ⚡";
            format-discharging = "{capacity}% 🔋";
            format-plugged = "{capacity}% 🔌";
            format-full = "CARREGADO";
            format-alt = "{time} {icon}";
            # "format-good": "", // An empty format will hide the module
            # "format-full": "",
            format-icons = ["" "" "" "" ""];
        };
        "battery#bat2" = {
            bat = "BAT2";
        };
        network = {
            # "interface": "*", // (Optional) To force the use of this interface
            format-wifi = "{essid} 🛜";
            format-ethernet = "🔗";
            tooltip-format = "{ifname} via {gwaddr} {ipaddr}/{cidr}";
            format-linked = "{ifname} (No IP)";
            format-disconnected = "Desconectado 🦖";
            format-alt = "{ifname}: {ipaddr}/{cidr}";
        };
        pulseaudio = {
            # "scroll-step": 10, // %, can be a float
            format = "{volume}%{icon} {format_source}";
            format-bluetooth = "{volume}% {icon} {format_source}";
            format-bluetooth-muted = "🔇 {icon} {format_source}";
            format-muted = "🔇 {format_source}";
            format-source = "{volume}% 🎙️";
            format-source-muted = "🔇";
            format-icons = {
                headphone = "🎧";
                hands-free = "";
                headset = "";
                phone = "📱";
                portable = "";
                car = "";
                default = ["" "" ""];
            };
            on-click = "pavucontrol";
        };
      };
    };
  };
}
