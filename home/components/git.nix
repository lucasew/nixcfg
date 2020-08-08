{ config, pkgs, ... }:
let
  cfg = import ../../config.nix;
in
{
  programs.git = {
    enable = true;
    userName = cfg.username;
    userEmail = cfg.email;
  };
}
