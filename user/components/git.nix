{ config, pkgs, ... }:
let
  commonConf = import ../common;
in
{
  programs.git = {
    enable = true;
    userName = commonConf.username;
    userEmail = commonConf.email;
  };
}
