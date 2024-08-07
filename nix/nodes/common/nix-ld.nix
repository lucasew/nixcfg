{config, pkgs, lib, ...}:

lib.mkIf config.programs.nix-ld.enable {
  environment.systemPackages = [ pkgs.nix-alien ];

  programs.nix-ld.libraries = with pkgs; [
    fuse
    xorg.libXi
    wayland
    alsa-lib
  ] ++ (appimageTools.defaultFhsEnvArgs.targetPkgs pkgs);

}
