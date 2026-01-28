{ pkgs, ... }:

{
  config = {
    environment.systemPackages = [ pkgs.workspaced ];

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
