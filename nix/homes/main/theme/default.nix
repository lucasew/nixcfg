{pkgs, ...}: {
  imports = [
    ../../../stylix.nix
  ];

  home.packages = with pkgs; [
    lxappearance
    custom.colorpipe
  ];

  gtk.enable = true;
}
