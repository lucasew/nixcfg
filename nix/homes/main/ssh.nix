{ global, lib, ... }:
let
  inherit (lib) concatStringsSep mapAttrsToList;
  mkMatchBlock = name: host: ''
    Host ${name}
      HostName ${host.tailscale_ip}
      ${if host ? port then "Port ${toString host.port}" else ""}
      ${if host ? user then "User ${host.user}" else "User ${global.username}"}
  '';
  sshConfig = concatStringsSep "\n" (mapAttrsToList mkMatchBlock global.hosts);
in
{
  programs.ssh = {
    enable = true;
    extraConfig = sshConfig;
  };
}
