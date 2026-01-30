{ ... }:

{
  config = {
    systemd.user.services.workspaced = {
      description = "Workspaced Daemon";
      wantedBy = [ "graphical-session.target" ];
      serviceConfig = {
        ExecStart = "/home/lucasew/.local/share/workspaced/bin/workspaced daemon";
        Restart = "on-failure";
      };
    };
  };
}
