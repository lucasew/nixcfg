{ config, pkgs, ... }:

{
  imports = [
    ./bash.nix
    ./compression.nix
    ./dconf.nix
    ./git.nix
    ./htop.nix
    ./mspaint.nix
    ./neovim
    ./pinball.nix
    ./rclone.nix
    ./rofi.nix
    ./spotify.nix
    ./stremio.nix
    ./tmux
    ./usb_tixati.nix
    ./vscode
  ];
}
