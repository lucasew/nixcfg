{ pkgs, ... }:
let cfg = import ../config.nix
{
  allowUnfree = cfg.allowUnfree;
}
