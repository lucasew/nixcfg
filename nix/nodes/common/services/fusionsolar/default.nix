{ config, lib, ... }:
let
  cfg = config.services.fusionsolar-reporter;
in
lib.mkIf cfg.enable {
  sops.secrets."fusionsolar" = {
    sopsFile = ../../../../../secrets/fusionsolar.env;
    format = "dotenv";
  };
  services.fusionsolar-reporter.environmentFile = "/var/run/secrets/fusionsolar";
}
