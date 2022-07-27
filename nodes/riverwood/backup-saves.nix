{config, pkgs, ...}:
{
  # TODO: abstract stuff to options
  systemd.user.services."backup-saves" = {
    enable = true;
    description = "Backup some paths to a Git repo";
    path = with pkgs; [ git openssh libnotify ];
    script = ''
      set -eu
      SSH_AUTH_SOCK=/run/user/$(id -u)/ssh-agent
      notify-send "Backup iniciado"
      ${pkgs.python3Packages.python.interpreter} ${../../bin/backup-savegames} || notify-send "Backup terminado com erros" && notify-send "Backup terminado com sucesso"
      
    '';
    startAt = "*-*-* *:00:00"; # hourly, i guess
    environment = {
      STATE_BACKUP_DIR = "~/WORKSPACE/SAVES";
      SAVEGAMES_FOLDERS = builtins.toJSON {
        "skyrim" = [
          "~/.config/The_Elder_Scrolls_V_Skyrim_Special_Edition_AppImage_01/users/lucasew/Documents/My Games/Skyrim Special Edition"
        ];
      };
    };
  };
}
