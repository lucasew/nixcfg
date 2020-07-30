{config, pkgs, ...}:

{
    programs.bash = {
        enable = true;
        initExtra = ''
            export EDITOR="nvim"
            export PS1="\u@\h \w \$?\\$ \[$(tput sgr0)\]"
        '';
    };
}