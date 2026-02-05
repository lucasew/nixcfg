{pkgs, ...}: {
  users = {
    users = {
      # vai incrementando
      /*
      w_you = {
        isNormalUser = true;
        extraGroups = [
          "work"
          "ssh"
        ];
        uid = 2003;
      };
      */
    };
    groups.work = {
      gid = 1999;
    };
  };
  security.pam.enableEcryptfs = true;
  boot.kernelModules = ["ecryptfs"];
  environment.systemPackages = with pkgs; [
    ecryptfs
    brave
  ];
  virtualisation.podman.enable = true;
}
