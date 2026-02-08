{
  pkgs,
  ...
}:
{
  home.packages = with pkgs; [
    lxappearance
    custom.colorpipe
  ];

  gtk.enable = true;
}
