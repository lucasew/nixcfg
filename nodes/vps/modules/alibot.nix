{pkgs, config, lib, ...}:
with import ../../../globalConfig.nix;
with lib;
let
  cfg = config.vps.alibot;
  alibot = pkgs.callPackage "${
    builtins.fetchGit {
      url = "ssh://git@github.com/lucasew/alibot";
      rev = "5bf5a883f7e600905280a9ea4a445f575e94a04d";
    }
  }/package.nix" {};
in
{
  options = {
    vps.alibot = {
      enable = mkEnableOption "Enable alibot";
      secretsDotenv = mkOption {
        type = types.str;
        description = "a dotenv file with a BOT_TOKEN variable";
        default = "${rootPath}/secrets/alibot.env";
      };
      stateStore = mkOption {
        type = types.str;
        description = "where to save the state file";
        default = "/persist/alibot.json";
      };
    };
  };
  config = mkIf cfg.enable {
    systemd = {
      services.alibot = {
        enable = true;
        serviceConfig = {
          Type = "simple";
          Restart = "always";
          ExecStart = "${pkgs.dotenv}/bin/dotenv '@${cfg.secretsDotenv}' -- ${alibot}/bin/alibot -d '${cfg.stateStore}'";
        };
        wantedBy = [
          "default.target"
        ];
      };
    };
  };
}
