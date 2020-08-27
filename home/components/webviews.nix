{pkgs, config, ...}:
let
  whatsapp = pkgs.stdenv.mkNativefier {
    name = "WhatsApp";
    url = "https://web.whatsapp.com";
    props = {};
  };
  remnote = pkgs.stdenv.mkNativefier {
    name = "RemNote";
    url = "https://www.remnote.io/";
    props = {};
  };
  notion = pkgs.stdenv.mkNativefier {
    name = "NotionSo";
    url = "https://notion.so";
    props = {};
  };
in
{
    home.packages = with pkgs; [
      whatsapp
      remnote
      notion
    ];
}
