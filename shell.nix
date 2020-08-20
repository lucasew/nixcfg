let
    global = import ./globalConfig.nix;
in global.defaultPkgs.mkShell {
    shellHook = ''
    ${global.setupScript}

    echo Ambiente carregado!
    '';
}