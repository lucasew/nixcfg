{ pkgs, ... }:

{
  users = {
    users = {
      # reservado uid 2000
      w_cilia = {
        isNormalUser = true;
        extraGroups = [
          "work"
          "ssh"
        ];
        uid = 2001;
      };
      w_logito = {
        isNormalUser = true;
        extraGroups = [
          "work"
          "ssh"
        ];
        uid = 2002;
      };
    };
    groups.work = {
      gid = 1999;
    };
  };
  security.pam.enableEcryptfs = true;
  boot.kernelModules = [ "ecryptfs" ];
  environment.systemPackages = with pkgs; [
    ecryptfs
    brave
    docker-compose
    ripgrep
    fd
    helix
    ruff
    uv
  ];
  virtualisation.podman.enable = true;
}
