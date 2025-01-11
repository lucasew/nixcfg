{ pkgs, ... }:
{
  boot = {
    kernel.sysctl = {
      "vm.swappiness" = 10;
    };
    tmp.cleanOnBoot = true;
  };
  services = {
    ananicy = {
      enable = true;
      package = pkgs.ananicy-cpp;
    };
  };
}
