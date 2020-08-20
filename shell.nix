let
    nixpkgs = builtins.fetchTarball {
        url = "https://github.com/NixOS/nixpkgs/archive/20.03.tar.gz";
        sha256 = "0182ys095dfx02vl2a20j1hz92dx3mfgz2a6fhn31bqlp1wa8hlq";
    };
    pkgs = import nixpkgs {};
    machine_name = "acer";
in 
pkgs.mkShell {
    shellHook = ''
    export NIXPKGS_ALLOW_UNFREE=1
    export NIXOS_CONFIG=$(pwd)/machine/${machine_name}
    NIX_PATH=$NIX_PATH:nixpkgs-overlays=$(pwd)/overlays:nixpkgs=${nixpkgs}:nixos-config=$NIXOS_CONFIG

    alias nixos-rebuild="sudo -E nixos-rebuild"
    alias nixos-install="sudo -E nixos-install --system $NIXOS_CONFIG"
    echo Ambiente carregado!

    '';
}