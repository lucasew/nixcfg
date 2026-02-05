{
  pkgs,
  global,
  lib,
  ...
}: let
  inherit (builtins) attrValues concatStringsSep;
  inherit (lib) flatten;
  hostAttrs = attrValues global.hosts;
  ips = flatten (
    map (
      h:
        (
          if (h ? tailscale_ip)
          then [h.tailscale_ip]
          else []
        )
        ++ (
          if (h ? zerotier_ip)
          then [h.zerotier_ip]
          else []
        )
    )
    hostAttrs
  );
  ipExprs = map (ip: "allow ${ip};") ips;
in {
  services.nginx = {
    package = pkgs.unstable.nginxMainline;
    appendHttpConfig = ''
      ${concatStringsSep "\n" ipExprs}
      allow 127.0.0.1;
      allow ::1;
      deny all;
    '';
  };
}
