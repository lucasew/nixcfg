path:
let
    pkgs = import <nixpkgs> {};
    kvs = if (builtins.pathExists path) then
        builtins.readDir path
    else
        abort("${path} not found");
    justKs = pkgs.lib.mapAttrsToList (k: v: k) kvs;
in justKs