{ lib, pkgs, ... }:

let
  sshfsArgs = lib.escapeShellArgs [
    "-f"
    "-o"
    "reconnect,ServerAliveInterval=15,ServerAliveCountMax=3,allow_other,default_permissions,cache=no"
  ];
in

{
  environment.systemPackages = [ pkgs.sshfs ];

  systemd.user.services = {
    "sshfs-TMP2" = {
      path = with pkgs; [ sshfs ];
      environment.SSH_AUTH_SOCK = "%t/ssh-agent";
      script = ''
        exec sshfs $(whoami)@whiterun:/home/$(whoami)/TMP2 /home/$(whoami)/TMP2 ${sshfsArgs}
      '';
      restartIfChanged = true;
    };

    "sshfs-WORKSPACE" = {
      path = with pkgs; [ sshfs ];
      environment.SSH_AUTH_SOCK = "%t/ssh-agent";
      script = ''
        exec sshfs $(whoami)@whiterun:/home/$(whoami)/WORKSPACE /home/$(whoami)/WORKSPACE ${sshfsArgs}
      '';
      restartIfChanged = true;
    };
  };
}
