let
  inherit (builtins) getFlake;
  flake = getFlake "${toString ./.}";
  inherit (flake.outputs) pkgs;
in builtins.attrValues {
  inherit (pkgs.x86_64-linux)
    stremio
    minecraft
    discord
  ;
  # inherit (pkgs.aarch64-linux.custom)
  #   neovim
  #   emacs
  # ;
  inherit (pkgs.x86_64-linux.python3Packages)
    scikitlearn
  ;
  inherit (pkgs.x86_64-linux.wineApps)
    wine7zip
    pinball
  ;
  polybar = pkgs.x86_64-linux.callPackage ./modules/polybar/customPolybar.nix {};
  inherit (flake.outputs.nixosConfigurations.x86_64-linux)
    acer-nix
    vps
  ;
  inherit (flake.outputs.homeConfigurations.x86_64-linux)
    main
  ;
}
