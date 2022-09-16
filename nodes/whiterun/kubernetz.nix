{config, pkgs, ...}:
let
  masterIp = "192.168.69.1";
  masterAPIServerPort = 6443;
in {
  services.kubernetes = {
    roles = [ "master" "node" ];
    masterAddress = masterIp;
    apiserverAddress = "http://${masterIp}:${toString masterAPIServerPort}";
    kubelet.extraOpts = "--fail-swap-on=false";
    apiserver = {
      securePort = masterAPIServerPort;
      advertiseAddress = masterIp;
    };
    addons.dns.enable = true;
  };
  environment.etc."cni/net.d".enable = false;
  environment.etc."cni/net.d/11-flannel.conf".source = "${config.environment.etc."cni/net.d".source}/11-flannel.conf";
  environment.systemPackages = with pkgs; [
    kompose
    kubectl
    kubernetes
  ];
}
