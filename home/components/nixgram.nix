{pkgs, config, ...}: 
let
  bin = pkgs.writeShellScriptBin "nixgram" ''
    # export PATH=${builtins.getEnv "PATH"}
    echo "Iniciando..."
    ${pkgs.dotenv}/bin/dotenv @${../../secrets/nixgram.env} -- ${pkgs.nixgram}/bin/nixgram
  '';
in
{
  config = {
    systemd.user.services.nixgram = {
      Unit = {
        Description = "Command bot for telegram";
        PartOf = ["graphical-session.target"];
      };
      Service = {
        Type = "exec";
        ExecStart = "${bin}/bin/nixgram";
        Restart = "on-failure";
      };
      Install = {
        WantedBy = [
          "default.target"
        ];
      };
    };
    home.packages = [pkgs.nixgram];
  };
}
