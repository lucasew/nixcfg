{pkgs, lib, ...}:
pkgs.wrapEmacs {
  imports = [
    ./fix-evil-org-mode-evil-redirect-digit-argument.nix
  ];
  evil = {
    enable = true;
    escesc = true;
    collection = true;
  };
  language-support = {
    nix.enable = true;
    markdown.enable = true;
  };
  org.enable = true;
  nogui = true;
  themes.selected = "wombat";
}
