{ self, config, pkgs, lib, ... }:
let
  inherit (lib) mkOption mkIf mkEnableOption types optionalString optional;
  ini = pkgs.formats.ini {};
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
      type = types.path;
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

      type = types.submodule {
        freeformType = ini.type;

        options = {
          general = {
            divider = mkOption {
              type = types.str;
              default = ",";
            };
          };

          search = mkOption {
            description = lib.mdDoc ''
              Search path related settings
            '';

            type = types.attrsOf (mkOption {
              type = types.listOf types.str;
              default = [];
              apply = builtins.concatStringsSep cfg.settings.general.divider;
            });

          };
        };
      };
    };
  };

  config = mkIf cfg.enable {
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
        paths = []
        ++ (optional cfg.enableGit pkgs.git);
        environment = {
          OUTPUT_DIR = cfg.outputDir;
          CONFIG_FILE = ini.generate "cloud-savegame-settings.ini" cfg.settings;
        };
        description = "Cloud savegame service";
        script = ''
          if [ -d "$OUTPUT_DIR" ]; then
            ${cfg.package}/bin/cloud-savegame -o "$OUTPUT_DIR" -c "$CONFIG_FILE" ${optionalString "-g" cfg.enableGit} ${optionalString "-v" cfg.enableVerbose}
          else
            echo "Output dir doesn't exist"
          fi
        '';
      };
    };
  };
}
