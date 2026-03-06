{
  config,
  pkgs,
  lib,
  ...
}:

lib.mkIf config.programs.sway.enable {
  programs.xss-lock = {
    enable = true;
    lockerCommand = lib.mkDefault "$HOME/.local/bin/lock-screen";
  };

  environment.systemPackages = [ pkgs.swaylock ];
}
