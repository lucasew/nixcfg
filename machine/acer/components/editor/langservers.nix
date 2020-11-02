{pkgs, ...}:
let
  pkgbin = import <dotfiles/lib/pkgbin.nix>;
in
{
  go = ["gopls"];
  rust = ["rls"];
  python = ["python-language-server"];
  nix = [(pkgbin "rnix-lsp")];
  c = [(pkgbin "ccls")];
  cpp = [(pkgbin "ccls")];
}
