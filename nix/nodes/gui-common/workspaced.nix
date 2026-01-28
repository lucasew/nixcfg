{ pkgs, ... }:

{
  config = {
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
      serviceConfig = {
        ExecStart = "${pkgs.workspaced}/bin/workspaced daemon";
        Restart = "on-failure";
      };
    };
  };
}
