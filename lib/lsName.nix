path:
let
    pkgs = import <nixpkgs> {};
    kvs = if (builtins.pathExists path) then
        builtins.readDir path
    else
        abort("${path} not found");
    justKs = pkgs.lib.mapAttrsToList (k: v: k) kvs;
    fn = k: path + "/${k}";
in map fn justKs
