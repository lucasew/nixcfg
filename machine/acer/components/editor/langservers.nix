{pkgs, ...}:
let
  pkgbin = name:
    let
      pkg = pkgs."${name}";
    in ["${pkg}/bin/${name}"];
in
{
  go = ["gopls"];
  rust = ["rls"];
  python = ["python-language-server"];
  nix = pkgbin "rnix-lsp";
}
