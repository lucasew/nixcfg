{...}:
{
  services.zerotierone = {
    enable = true;
    port = 6969;
    joinNetworks = [
      "e5cd7a9e1c857f07"
    ];
  };
  networking.firewall.trustedInterfaces = [ "ztppi77yi3" ];
  networking.extraHosts = ''
    192.168.69.1 controlplane.local
    192.168.69.1 whiterun.local
    192.168.69.2 riverwood.local
    192.168.69.4 xiaomi.local
  '';
}
