{...}: {
  services.dnsmasq = {
    enable = true;
    servers = [ "8.8.8.8" "8.8.4.4" ];
    extraConfig = ''
domain-needed
bogus-priv
addn-hosts=/etc/extraHosts
    '';
  };
}
