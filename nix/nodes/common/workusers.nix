{ pkgs, ...}:

{
  users = {
    users = {
      w_dr = {
        isNormalUser = true;
        extraGroups = [
          "work"
        ];
        uid = 2000;
      };
      w_cilia = {
        isNormalUser = true;
        extraGroups = [
          "work"
        ];
        uid = 2001;
      };
    };
    groups.work = {
      gid = 1999;
    };
  };
  security.pam.enableEcryptfs = true;
  boot.kernelModules = [ "ecryptfs" ];
  environment.systemPackages = [ pkgs.ecryptfs pkgs.chromium ];
}
