{pkgs, ...}: 
let
  dotenv = import ../dotenv/package.nix;
  nixgram = import ./package.nix;
  bin = pkgs.writeShellScriptBin "nixgram" ''
    # export PATH=${builtins.getEnv "PATH"}
    echo "Iniciando..."
    ${dotenv}/bin/dotenv @${../../secrets/nixgram.env} -- ${nixgram}/bin/nixgram
  '';
  systemdUserService = import <dotfiles/lib/systemdUserService.nix>;
in systemdUserService {
  description = "Command bot for telegram";
  command = "${bin}/bin/nixgram";
}
