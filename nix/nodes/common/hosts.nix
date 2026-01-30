{
  global,
  config,
  lib,
  ...
}:
let
  node = global.hosts.${config.networking.hostName}.tailscale_ip or null;

  nginxDomains = builtins.attrNames config.services.nginx.virtualHosts;
  baseDomain = "${config.networking.hostName}.${config.networking.domain}";
  allMySubdomains = lib.flatten [
    nginxDomains
    baseDomain
  ];

  tinydnsLines = map (item: "+${item}:${node}:${toString ttl}") allMySubdomains;
  tinydnsData =
    if node != null then (builtins.concatStringsSep "\n" (lib.unique tinydnsLines)) else "";

  ttl = 30;
in
lib.mkIf (node != null) {
  services.tinydns = {
    data = ''
      .${baseDomain}:${node}:ns:${toString ttl}
      ${tinydnsData}
    '';
  };
}
