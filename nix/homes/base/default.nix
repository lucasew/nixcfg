{
  global,
  pkgs,
  config,
  self,
  lib,
  bumpkin,
  ...
}:
let
  inherit (self) inputs outputs;
  inherit (lib) mkDefault;
  inherit (global) environmentShell;
in
{
  imports = [
    "${self.inputs.nixgram}/hmModule.nix"
    "${self.inputs.redial-proxy}/hmModule.nix"
    "${self.inputs.borderless-browser}/home-manager.nix"
  ];

  home.packages = with pkgs; [
    yt-dlp # video downloader
    file # what file is it?
    neofetch # system info, arch linux friendly
    comma # like nix-shell but more convenient
    fzf # file finder and terminal based dmenu
    ffmpeg # video converter
    send2kindle
    home-manager
  ];

  home.stateVersion = mkDefault "22.11";
  home.enableNixpkgsReleaseCheck = false;

  #home.file.".dotfilerc".text = ''
  #  #!/usr/bin/env bash
  #  ${environmentShell}
  #'';
  programs = {
    tmux.enable = true;
    git = {
      enable = true;
      userName = global.username;
      userEmail = global.email;
      package = mkDefault pkgs.gitMinimal;
    };
  };
}
