{ config, pkgs, ... }:
{
  home.packages = with pkgs; [
    boringssl
  ];
  programs.vscode = {
    enable = true;
    package = pkgs.latest.vscode;
    extensions = (import ./extensions.nix) pkgs;
    userSettings = import ./userSettings.nix;
  };
}
