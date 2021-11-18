{pkgs, config, lib, ...}:
let
  inherit (lib) mkIf;
in mkIf (config.evil.enable && config.org.enable) {
  warnings = [
    "there is a bug on https://github.com/Somelauw/evil-org-mode and the fix should be available soon. See https://github.com/Somelauw/evil-org-mode/issues/93#issuecomment-950306532"
  ];
  initEl.pre = ''
    (fset 'evil-redirect-digit-argument 'ignore) ;; before evil-org loaded
  '';
}
