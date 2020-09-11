{ pkgs, ...}:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
  generator = import ./gen.nix globalConfig;
in
{
  home.file.".dotfilerc".text = ''
    #!/usr/bin/env bash
    ${generator}
  '';
}
