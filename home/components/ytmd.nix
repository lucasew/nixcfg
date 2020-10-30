{config, pkgs, lib, ...}:
with lib;
let 
    adskipper = 
        pkgs.writeShellScriptBin "ytmd-adskip" ''
            PLAYERCTL=${pkgs.playerctl}/bin/playerctl
            echo Executando...
            function handle {
                echo $1
                if [[ $1 =~ 'https://music.youtube.com/' ]]; then
                    echo Pulando ad...
                    $PLAYERCTL next -p youtubemusic
                fi
            }
            $PLAYERCTL metadata -p youtubemusic --format "{{mpris:artUrl}}" -F 2> /dev/null \
            | while read line; do \
                handle $line; \
            done
            '';
in {
    config = {
        systemd.user.services.ytmd-adblock = {
            Unit = {
                Description = "Youtube music ad skipper";
                PartOf = [ "graphical-session.target" ];
            };
            Service = {
                Type = "exec";
                ExecStart = "${adskipper}/bin/ytmd-adskip";
                Restart = "on-failure";
            };
            Install = {
              WantedBy = [
                "default.target"
              ];
            };
        };
        home.packages = [
          pkgs.ytmd
        ];
    };
}
