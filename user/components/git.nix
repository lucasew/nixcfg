{ config, pkgs, ... }:
{
  programs.git = {
    enable = true;
    userName = "lucasew";
    userEmail = "lucas59356@gmail.com";
  };
}
