{config, pkgs, lib, ...}:
let
  cfg = config.services.simple-dashboardd;
  module = pkgs.buildGoModule {
    pname = "simple-dashboard";
    version = "unstable-2022-07-15";
    src = pkgs.fetchFromGitHub {
      owner = "lucasew";
      repo = "simple-dashboard";
      rev = "0f5483a785d95ff53c18860ae737b675265d831c";
      sha256 = "sha256-r74QSiBj2N1xU9o7k/sO40ctKThbox+NvN8mJ8ho+OU=";
    };
    vendorSha256 = "sha256-a6iSGI+PJxIqF2WDp86SCR7Q2+pYf2kn0d7jKPScCyg=";
    postInstall = ''
      mkdir $out/share/simple-dashboard -p
      cp $src/*.ini* $out/share/simple-dashboard
    '';
    meta = with lib; {
      description = "Simple web-based dashboard to watch with your old tablet";
      homepage = "https://github.com/lucasew/simple-dashboard";
      license = licenses.mit;
      maintainers = with maintainers; [ lucasew ];
    };
  };
in {
  options = with lib; {
    services.simple-dashboardd = {
      enable = mkEnableOption "Webapp to show system usage";
      config = mkOption {
        description = "Config string to be used by simple-dashboard";
        type = types.lines;
        default = builtins.readFile "${module}/share/simple-dashboard/config.ini.example";
      };
      port = mkOption {
        description = "Port to listen";
        default = 42069;
        type = types.port;
      };
    };
  };
  config = with lib; mkIf cfg.enable {
    systemd.services.simple-dashboardd = {
      enable = true;
      path = [ module ];
      script = ''
        ${module}/bin/simple-dashboardd -c ${builtins.toFile "simple-dashboard.cfg" cfg.config} -p ${builtins.toString cfg.port}
      '';
      restartIfChanged = true;
    };
  };
}
