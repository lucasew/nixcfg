{ config, pkgs, ... }:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
in
{
  programs.bash = {
    enable = true;
    initExtra = ''
      export EDITOR="nvim"
      export PS1="\u@\h \w \$?\\$ \[$(tput sgr0)\]"
      export DOTFILES=${globalConfig.dotfileRootPath}
    '';
  };
}
