{
  pkgs,
  lib,
  ...
}: let
  inherit (lib) mkDefault;
in {
  home.packages = with pkgs; [
    neofetch # system info, arch linux friendly
    home-manager
  ];

  home.stateVersion = mkDefault "22.11";
  home.enableNixpkgsReleaseCheck = false;
}
