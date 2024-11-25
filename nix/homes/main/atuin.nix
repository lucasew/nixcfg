{ lib, config, ... }:

lib.mkIf config.programs.atuin.enable {
  programs.atuin = {
    enableBashIntegration = true;
  };
}

