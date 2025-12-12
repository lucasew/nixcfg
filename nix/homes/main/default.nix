{
  global,
  pkgs,
  lib,
  ...
}:

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
    ./qutebrowser.nix
    ./zen-browser.nix
    ./mise.nix
  ];

  stylix.enable = true;

  borderless-browser.chromium = lib.getExe pkgs.vivaldi;

  # programs.ghostty.enable = true;

  programs.atuin.enable = true;

  programs.discord.enable = true;
  # programs.zen-browser.enable = true; # Handled unconditionally in zen-browser.nix
  programs.vscode.enable = true;
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
}
