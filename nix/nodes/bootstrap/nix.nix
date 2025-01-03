{
  global,
  pkgs,
  lib,
  ...
}:
let
  inherit (global) username;
in
{
  nix = {
    package = pkgs.lix;
    gc = {
      automatic = true;
      dates = "weekly";
      options = "--delete-older-than 3d --max-freed ${toString (10 * 1024 * 1024 * 1024)}";
    };

    settings = {
      min-free = lib.mkDefault (5 * 1024 * 1024 * 1024);
      max-free = lib.mkDefault (100 * 1024 * 1024 * 1024);
      trusted-users = [
        username
        "@wheel"
      ];
      experimental-features = [
        "nix-command"
        "flakes"
        "impure-derivations"
        "ca-derivations"
      ];
    };
  };
}
