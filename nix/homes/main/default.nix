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
  ];

  borderless-browser.chromium = lib.getExe pkgs.brave;

  # programs.ghostty.enable = true;

  programs.atuin.enable = true;

  programs.zen-browser.enable = true;
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
    blender-bin.blender_3_6
    brave
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

    # dev
    conda
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

  services.redial_proxy.enable = true;

  programs = {
    # adskipped-spotify.enable = true;
    jq.enable = true;
    obs-studio = {
      package = pkgs.obs-studio;
      enable = true;
    };
    htop = {
      enable = true;
      settings = {
        hideThreads = true;
        treeView = true;
      };
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
