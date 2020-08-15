{ config, pkgs, ... }:

{
  imports = [
    ./bash.nix
    ./dconf.nix
    ./git.nix
    ./htop.nix
    ./neovim
    ./tmux
    ./vscode
    ./rofi.nix
    ./compression.nix
    ./rclone.nix
    ./spotify.nix
    ./stremio.nix
    ./pinball.nix
    ./usb_tixati.nix
  ];
}
