let
  pkgs = import <dotfiles/pkgs.nix>;
in
  {
    name, 
    url, 
    electron ? pkgs.electron_8, 
    props ? {}
  }:
let
  nativefier = pkgs.nodePackages.nativefier;
  processedProps = {
    name = name;
    targetUrl = url;
    badge = false;
    showMenuBar = false;
  } // props;
  writtenProps = pkgs.writeText "nativefier.json" (builtins.toJSON processedProps);
  drv = pkgs.stdenv.mkDerivation {
    name = name;
    dontUnpack = true;
    installPhase = ''
      export OUT_DIR=$out/share/${name}-nativefier
      mkdir $OUT_DIR -p
      cp -r ${nativefier}/lib/node_modules/nativefier/app/* $OUT_DIR
      rm $OUT_DIR/nativefier.json
      cat ${writtenProps} > $OUT_DIR/nativefier.json
    '';
  };
  binary = pkgs.writeShellScriptBin name ''
    ${electron}/bin/electron ${drv}/share/${name}-nativefier/lib/main.js
  '';
in binary
