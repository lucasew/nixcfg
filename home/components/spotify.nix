{config, pkgs, lib, ...}:
with lib;
let 
    cfg = config.programs.spotify;
    adskipper = 
        pkgs.writeShellScriptBin "spotify-adskip" ''
            PLAYERCTL=${pkgs.playerctl}/bin/playerctl
            echo Executando...
            function handle {
                echo $1
                if [[ $1 =~ ^spotify:ad:.* ]]; then
                    echo Pulando ad...
                    $PLAYERCTL next -p spotify
                fi
            }
            $PLAYERCTL metadata -p spotify --format "{{mpris:trackid}}" -F 2> /dev/null \
            | while read line; do \
                handle $line; \
            done
            '';
    adskipperBinary = "${adskipper}/bin/spotify-adskip";
in {
    options = {
        programs.spotify = {
            enable = mkEnableOption "Enable adblock for spotify";
        };
    };
    config = {
        systemd.user.services.spotify-adblock = {
            Unit = {
                Description = "Spotify ad skipper";
            };
            Service = {
                Type = "exec";
                ExecStart = adskipperBinary;
                Restart = "on-failure";
            };
        };
        home.packages = [pkgs.spotify];
    };
}