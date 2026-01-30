{ pkgs, ... }:

{
  config = {
    environment.systemPackages = [ pkgs.workspaced ];

    systemd.user.services.workspaced = {
      description = "Workspaced Daemon";
      wantedBy = [ "graphical-session.target" ];
      path = with pkgs; [ 
        workspaced 
        mise
      ];
      serviceConfig = {
        ExecStart = "workspaced daemon";
        Restart = "on-failure";
      };
    };
  };
}
