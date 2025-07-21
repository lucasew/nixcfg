{
  config,
  lib,
  self,
  ...
}:

let
  cfg = config.services.phpelo;
in

{
  imports = [
    "${self.inputs.phpelo}/nixos.nix"
  ];
  config = lib.mkIf cfg.enable {

    services.ts-proxy.hosts."phpelo-${config.networking.hostName}" = {
      enableTLS = true;
      enableFunnel = true;
      network = "unix";
      proxies = [ "phpelo.socket" ];
      address = cfg.socket;
    };
  };
}
