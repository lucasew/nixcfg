{ self, config, pkgs, lib, ... }:
let
  inherit (lib) mkOption mkIf mkEnableOption types optionalString optional;
  ini = pkgs.formats.ini {
    listToValue = value: builtins.concatStringsSep config.options.services.cloud-savegame.settings.general.divider;
  };
  cloud-savegame = pkgs.callPackage "${self.inputs.cloud-savegame}/package.nix" {};
  cfg = config.services.cloud-savegame;
in {
  options.services.cloud-savegame = {
    enable = mkEnableOption "Cloud savegame";

    calendar = mkOption {
      type = types.str;
      default = "00:00:01";
      description = lib.mdDoc ''
        When to run Cloud Savegame in systemd timer format
      '';
    };

    package = mkOption {
      type = types.package;
      default = cloud-savegame;
      defaultText = "* from the flake input *";
      description = lib.mdDoc ''
        Cloud Savegame package to use
      '';
    };

    enableVerbose = mkEnableOption "Enable verbose output for Cloud savegame";

    enableGit = mkEnableOption "Enable git interactions for Cloud savegame";

    outputDir = mkOption {
      type = types.str;
      default = "~/SavedGames";
      description = ''
        Where the savegame repo will be stored
      '';
    };

    settings = mkOption {
      description = lib.mdDoc ''
        Cloud savegame settings

        These are converted to the ini file
      '';

      type = ini.type;

      example = builtins.fromTOML "${self.inputs.cloud-savegame}/demo.cfg";

      default = {
        general.divider = ",";
        search.paths="~";
        flatout-2.installdir= ["~/.local/share/Steam/steamapps/common/FlatOut2" ];
      };
      # TODO: Convert all lists to strings divided by the divider using 'builtins.concatStringsSep cfg.settings.general.divider'
    };
  };

  config = mkIf cfg.enable {
    services.cloud-savegame.settings = {
      general.divider = ",";
      search.paths="~";
      flatout-2.installdir= ["~/.local/share/Steam/steamapps/common/FlatOut2" ];
    };

    environment.etc."cloud-savegame-settings.ini".source = 
      let
        atom = val: if ((builtins.typeOf val) == "list") then
          (builtins.concatStringsSep (cfg.settings.general.divider or ",") (map (toString) val))
        else if ((builtins.typeOf val) == "set") then
          (builtins.mapAttrs (k: v: atom v) val)
        else (toString val);
      in ini.generate "cloud-savegame-settings.ini" (atom cfg.settings);

    systemd.user = {
      timers.cloud-savegame = {
        description = "Cloud savegame timer";
        wantedBy = [ "timers.target" ];
        timerConfig = {
          OnCalendar = cfg.calendar;
          AccuracySec = "30m";
          Unit = "cloud-savegame.service";
        };
      };

      services.cloud-savegame = {
        enable = true;
        path = []
        ++ (optional cfg.enableGit pkgs.git)
        ++ (optional cfg.enableGit pkgs.openssh);
        environment = {
          OUTPUT_DIR = cfg.outputDir;
          CONFIG_FILE = "/etc/cloud-savegame-settings.ini";
        };
        description = "Cloud savegame service";
        script = ''
          export SSH_AUTH_SOCK=/run/user/$(id -u)/ssh-agent
          tilde=~
          OUTPUT_DIR="$(echo "$OUTPUT_DIR" | sed "s;~;$tilde;")"
          if [ -d "$OUTPUT_DIR" ]; then
            ${cfg.package}/bin/cloud-savegame -o "$OUTPUT_DIR" -c "$CONFIG_FILE" ${optionalString cfg.enableGit "-g"} ${optionalString cfg.enableVerbose "-v"}
          else
            echo "Output dir '$OUTPUT_DIR' doesn't exist"
          fi
        '';
      };
    };
  };
}
