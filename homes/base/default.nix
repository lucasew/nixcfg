{ global, pkgs, config, self, lib, ...}:
let
  inherit (self) inputs outputs;
  environmentShell = outputs.environmentShell.x86_64-linux;
  inherit (lib) mkDefault;
in
{
  home.packages = with pkgs; [
    youtube-dl # video downloader
    file # what file is it?
    neofetch # system info, arch linux friendly
    comma # like nix-shell but more convenient
    fzf # file finder and terminal based dmenu
  ];

  home.stateVersion = mkDefault "20.03";

  home.file.".dotfilerc".text = ''
    #!/usr/bin/env bash
    ${environmentShell}
  '';

}