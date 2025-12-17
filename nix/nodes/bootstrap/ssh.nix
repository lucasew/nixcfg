{ config, ... }:
{
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = false;
    };
  };

  programs.ssh.extraConfig = ''
    Include ${config.sops.secrets.ssh-hosts.path}
  '';
  programs.mosh.enable = true;

  users.groups.ssh = { };

  sops.secrets.ssh-hosts = {
    sopsFile = ../../../secrets/ssh-hosts;
    owner = config.users.users.root.name;
    group = config.users.groups.ssh.name;
    mode = "0440";
    format = "binary";
  };
}
