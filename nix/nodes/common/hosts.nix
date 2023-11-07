{ global, config, lib, ... }:
let
  node = global.nodeIps.${config.networking.hostName}.ts or null;

  nginxDomains = builtins.attrNames config.services.nginx.virtualHosts;
  baseDomain = "${config.networking.hostName}.${config.networking.domain}";
  allMySubdomains = lib.flatten [ nginxDomains baseDomain ];

  tinydnsLines = map (item: ".${item}:${node}") allMySubdomains;
in {
  services.tinydns = {
    enable = lib.mkDefault node != null;
    data = if node != null then (builtins.concatStringsSep "\n"  tinydnsLines) else "";
  };
}
