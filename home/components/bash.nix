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
      source ~/.dotfilerc
      alias la="ls -a"
      alias ncdu="rclone ncdu . 2> /dev/null"
    '';
  };
}
