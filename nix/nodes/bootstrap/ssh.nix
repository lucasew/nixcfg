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
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = true;
    };
  };

  programs.ssh = {
    extraConfig = sshConfig;
  };

  programs.mosh.enable = true;
}
