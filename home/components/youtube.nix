{config, pkgs, lib, ...}:
let
  systemdUserService = import <dotfiles/lib/systemdUserService.nix>;
in
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
          if [[ $($PLAYERCTL metadata -p youtubemusic mpris:length) =~ 30000000 ]]; then
            echo Ad com timeout de 5s
            sleep 5
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
      systemd.user.services.ytmd-adblock = systemdUserService {
          description = "Youtube music ad skipper";
          command = "${adskipper}/bin/ytmd-adskip";
        };
        home.packages = [
          pkgs.ytmd
        ];
    };
}
