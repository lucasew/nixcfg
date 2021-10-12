{ pkgs ? import <nixpkgs> {} }:
let
in
pkgs.vscode-utils.vscodeEnv
{
  nixExtensions = import ./extensions.nix {inherit pkgs;};
  settings = import ./userSettings.nix {inherit pkgs;};
  mutableExtensionsFile = "/tmp/code/state";
  user-data-dir = "/tmp/code/user-data-dir";
}
