{pkgs, config, ...}: 
let
  bin = pkgs.writeShellScriptBin "nixgram" ''
    # export PATH=${builtins.getEnv "PATH"}
    echo "Iniciando..."
    ${pkgs.dotenv}/bin/dotenv @${../../secrets/nixgram.env} -- ${pkgs.nixgram}/bin/nixgram
  '';
  systemdUserService = import <dotfiles/lib/systemdUserService.nix>;
in
{
  config = {
    systemd.user.services.nixgram = systemdUserService {
        description = "Command bot for telegram";
        command = "${bin}/bin/nixgram";
    };
    home.packages = [pkgs.nixgram];
  };
}
