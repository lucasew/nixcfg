{pkgs, config, ...}:
let
  globalConfig = import <dotfiles/globalConfig.nix>;
  fetch = url: builtins.fetchurl {url = url;};
in
let
  whatsapp = pkgs.stdenv.mkNativefier {
    name = "WhatsApp";
    url = "https://web.whatsapp.com";
    electron = pkgs.electron_9;
    icon = fetch "https://raw.githubusercontent.com/jiahaog/nativefier-icons/gh-pages/files/whatsapp.png";
    props = {
      userAgent = "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36";
      singleInstance = true;
      # tray = true;
    };
  };
  remnote = pkgs.stdenv.mkNativefier {
    name = "RemNote";
    electron = pkgs.electron_9;
    url = "https://www.remnote.io/";
    icon = fetch "https://www.remnote.io/favicon.ico";
  };
  notion = pkgs.stdenv.mkNativefier {
    name = "NotionSo";
    url = "https://notion.so";
    icon = fetch "https://logos-download.com/wp-content/uploads/2019/06/Notion_App_Logo.png";
  };
  duolingo = pkgs.stdenv.mkNativefier {
    name = "Duolingo";
    url = "https://duolingo.com";
    icon = fetch "https://logos-download.com/wp-content/uploads/2016/10/Duolingo_logo_owl.png";
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
