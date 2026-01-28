{ config, lib, pkgs, ... }:

{
  config = lib.mkIf (config.programs.sway.enable || config.services.xserver.windowManager.i3.enable) {
    environment.systemPackages = [ pkgs.workspaced ];

    systemd.user.sockets.workspaced = {
      description = "Workspaced Socket";
      wantedBy = [ "sockets.target" ];
      listenStreams = [ "%t/workspaced.sock" ];
    };

    systemd.user.services.workspaced = {
      description = "Workspaced Daemon";
      wantedBy = [ "graphical-session.target" ];
      requires = [ "workspaced.socket" ];
      restartTriggers = [ pkgs.workspaced ];
      serviceConfig = {
        ExecStart = "${pkgs.workspaced}/bin/workspaced daemon";
        Restart = "on-failure";
      };
    };
  };
}
