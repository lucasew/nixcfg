{
  pkgs,
  ...
}:
{
  home.packages = with pkgs; [
    lxappearance
    custom.colorpipe
  ];

  # Desabilita gerenciamento de mimeapps e gtk - gerenciado pelo workspaced
  xdg.mimeApps.enable = false;
}
