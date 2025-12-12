{
  config,
  pkgs,
  lib,
  ...
}:

{
  home.activation = {
    setup-zen-browser = lib.hm.dag.entryAfter [ "writeBoundary" ] ''
      PATH+=":"~".local/share/flatpak/exports/bin:/var/lib/flatpak/exports/bin:/run/current-system/sw/bin"
      zenBin=io.github.zen_browser.zen
      if /run/current-system/sw/bin/sdw source_me has_binary $zenBin; then
        run xdg-settings set default-web-browser $zenBin.desktop
      else
        echo WARNING: zen browser is not installed: flatpak install $zenBin
      fi
    '';
  };
}
