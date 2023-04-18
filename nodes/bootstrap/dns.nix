{ config, ... }:
let
  whiterun = "100.85.38.19";
  riverwood = "100.107.51.95";
in {
  # DNS based on MagicDNS
  services.dnsmasq.extraConfig = ''
address=/whiterun.${config.networking.domain}/${whiterun}
address=/riverwood.${config.networking.domain}/${riverwood}
  '';

}
