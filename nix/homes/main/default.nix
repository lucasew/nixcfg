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
    ./ghostty.nix
    ./dconf.nix
    ./theme
    ./qutebrowser.nix
    ./zen-browser.nix
    ./ssh.nix
  ];

  home.activation.restart-workspaced = lib.hm.dag.entryAfter ["writeBoundary"] ''
    dotfilesFolder=
    if [ -d ~/.dotfiles ]; then
      dotfilesFolder=~/.dotfiles
    elif [ -d /home/lucasew/.dotfiles ]; then
      dotfilesFolder=/home/lucasew/.dotfiles
    elif [ -d /etc/.dotfiles ]; then
      dotfilesFolder=/etc/.dotfiles
    fi
    if [ -n "$dotfilesFolder" ]; then
      mkdir -p ~/.local/share/workspaced/bin
      (cd "$dotfilesFolder/nix/pkgs/workspaced" && "$dotfilesFolder/bin/shim/mise" exec -- go build -o ~/.local/share/workspaced/bin/workspaced ./cmd/workspaced)
    fi
    $DRY_RUN_CMD systemctl --user restart workspaced.service || true
  '';

  stylix.enable = true;

  programs.atuin.enable = true;

  programs.man.enable = true;

  home = {
    homeDirectory = /home/lucasew;
    inherit (global) username;
  };

  home.packages = with pkgs; [
    mission-center
    cached-nix-shell
    feh
    fortune
    graphviz
    github-cli
    google-cloud-sdk
    htop
    libnotify
    ncdu
    nix-prefetch-scripts
    nix-output-monitor
    pkg
    remmina
    sqlite
    sshpass
    zenity

    # media
    nbr.wine-apps._7zip
    xxd

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

  programs = {
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
