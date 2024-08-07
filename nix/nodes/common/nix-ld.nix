{config, pkgs, lib, ...}:

lib.mkIf config.programs.nix-ld.enable {
  environment.systemPackages = [ pkgs.nix-alien ];
}
