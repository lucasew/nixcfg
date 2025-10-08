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
    "${self.inputs.borderless-browser}/home-manager.nix"
  ];

  programs.bash.bashrcExtra = ''
    if command -v sdw > /dev/null 2> /dev/null && [ -f "$(sdw d root || echo "/nhaa")/bin/source_me" ]; then
      source $(sdw d root)/bin/source_me
    fi
  '';

  home.packages = with pkgs; [
    file # what file is it?
    neofetch # system info, arch linux friendly
    comma # like nix-shell but more convenient
    fzf # file finder and terminal based dmenu
    send2kindle
    home-manager
  ];

  home.stateVersion = mkDefault "22.11";
  home.enableNixpkgsReleaseCheck = false;

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
