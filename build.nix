with import ./default.nix;
builtins.attrValues {
  inherit (pkgs)
    stremio
    minecraft
    discord
  ;
  inherit (pkgs.python3Packages)
    scikitlearn
  ;
  inherit (pkgs.wineApps)
    wine7zip
    pinball
  ;
  polybar = pkgs.callPackage ./modules/polybar/customPolybar.nix {};
  inherit (nixosConfigurations)
    acer-nix
    vps
  ;
  inherit (homeConfigurations)
    main
  ;
}
