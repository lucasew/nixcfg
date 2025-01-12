{ config, ... }:
{
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = true;
    };
  };

  programs.ssh.extraConfig = ''
    Include ${config.sops.secrets.ssh-hosts.path}
  '';
  programs.mosh.enable = true;

  sops.secrets.ssh-hosts = {
    sopsFile = ../../../secrets/ssh-hosts;
    owner = config.users.users.lucasew.name;
    group = config.users.users.lucasew.group;
    format = "binary";
  };
}
