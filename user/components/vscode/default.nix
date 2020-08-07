{ config, pkgs, ... }:
{
  imports = [
    # ./extensions.nix
    # ./config.nix
  ];
  programs.vscode = {
    enable = true;
    package = pkgs.vscode;
    extensions = (import ./extensions.nix) pkgs;
    userSettings = import ./userSettings.nix;
  };
}
