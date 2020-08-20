{ config, pkgs, ... }:
with pkgs.globalConfig;
{
  programs.git = {
    enable = true;
    userName = username;
    userEmail = email;
  };
}
