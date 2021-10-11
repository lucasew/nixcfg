{ pkgs ? import <nixpkgs> {} }:
pkgs.vscode-utils.vscodeEnv
{
  vscode = pkgs.vscode-with-extensions.override {
    vscodeExtensions = import ./extensions.nix {inherit pkgs;};
  };
  settings = import ./userSettings.nix {inherit pkgs;};
  mutableExtensionsFile = "/tmp/code/state";
  user-data-dir = "/tmp/code/user-data-dir";
}
