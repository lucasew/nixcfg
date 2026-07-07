{
  lib,
  self,
  pkgs,
  config,
  ...
}:
let
  inherit (lib) mkIf;
in
{
  imports = [ "${self.inputs.telegram-sendmail}/nixos-module.nix" ];

  config = mkIf config.services.telegram-sendmail.enable {
    services.telegram-sendmail.credentialFile = config.sops.secrets.telegram-sendmail.path;

    sops.secrets.telegram-sendmail = {
      sopsFile = ../../../../secrets/telegram_sendmail.env;
      owner = config.users.users.telegram_sendmail.name;
      group = config.users.users.telegram_sendmail.group;
      format = "dotenv";
    };
  };
}
