let
    global = import ./globalConfig.nix;
in global.defaultPkgs.mkShell {
    shellHook = ''
    export DOTFILES=$(pwd)
    ${global.setupScript}

    echo Ambiente carregado!
    '';
}
