{ ... }:
{
  boot = {
    kernel.sysctl = {
      "vm.swappiness" = 10;
    };
    tmp.cleanOnBoot = true;
  };
}
