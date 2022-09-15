{config, lib, pkgs, ...}:
let
  inherit (lib) types mkEnableOption mkOption mkIf;
  cfg = config.services.unstore;
  inherit (builtins) concatStringsSep replaceStrings;
in {
  options.services.unstore = {
    enable = mkEnableOption "unstore: scheduled delete of nix-store paths that contain a file pattern";
    paths = mkOption {
      description = "Path patterns to remove";
      type = types.listOf lib.types.str;
      default = [ "job_runner.ipynb" "flake.nix" ];
    };
    startAt = mkOption {
      description = "When to run the service";
      type = types.str;
      default = "*-*-* *:00:00";
    };
  };
  config = mkIf cfg.enable {
    systemd.services.unstore = {
      inherit (cfg) enable startAt;
      path = [ pkgs.nix ];
      description = "delete paths that contain a file pattern in the nix-store";
      script = let
        paths = builtins.concatStringsSep " " (map (path: "/nix/store/*/${builtins.replaceStrings [" "] ["\\ "] path}") cfg.paths);
      in ''
        for p in ${paths} ; do
          echo "Removing '$p'"
          nix-store --delete "$p" || true
        done
      '';
    };
  };
}
