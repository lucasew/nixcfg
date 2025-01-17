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
    owner = config.users.users.root.name;
    group = config.users.groups.users.name;
    mode = "0440";
    format = "binary";
  };
}
