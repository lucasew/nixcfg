{ pkgs, ...}:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
in
{
  home.file.".dotfilerc".text = ''
    #!/usr/bin/env bash
    alias nixos-rebuild="sudo -E nixos-rebuild --flake '${builtins.toString ./.}#acer-nix'"
  '';
}
