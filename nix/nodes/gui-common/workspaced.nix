{ lib, pkgs, ... }:

{
  config = {
    environment.systemPackages = [ pkgs.workspaced ];

    systemd.user.services.workspaced = {
      description = "Workspaced Daemon";
      wantedBy = [ "graphical-session.target" ];
      serviceConfig = {
        ExecStart = lib.getExe pkgs.workspaced;
        Restart = "on-failure";
      };
    };
  };
}
