{ config, pkgs, ... }:

{
	imports = [
		./neovim
		./tmux
		./bash.nix
		./dconf.nix
		./git.nix
		./vscode.nix
	];
}