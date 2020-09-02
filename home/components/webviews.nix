{pkgs, config, ...}:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
in
let
  whatsapp = pkgs.stdenv.mkNativefier {
    name = "WhatsApp";
    url = "https://web.whatsapp.com";
    electron = pkgs.latest.electron_9;
    props = {
      userAgent = "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36";
      singleInstance = true;
      # tray = true;
    };
  };
  remnote = pkgs.stdenv.mkNativefier {
    name = "RemNote";
    electron = pkgs.latest.electron_9;
    url = "https://www.remnote.io/";
  };
  notion = pkgs.stdenv.mkNativefier {
    name = "NotionSo";
    url = "https://notion.so";
  };
  duolingo = pkgs.stdenv.mkNativefier {
    name = "Duolingo";
    url = "https://duolingo.com";
  };
in
{
    home.packages = with pkgs; [
      whatsapp
      remnote
      notion
      duolingo
    ];
}
