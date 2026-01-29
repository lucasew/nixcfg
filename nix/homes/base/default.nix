{
  pkgs,
  lib,
  ...
}:
let
  inherit (lib) mkDefault;
in
{
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
    home-manager
    gitMinimal
  ];

  home.stateVersion = mkDefault "22.11";
  home.enableNixpkgsReleaseCheck = false;
}
