{...}:
let
  global = import <dotfiles/globalConfig.nix>;
in
{
  home.file.".nix-channels".text = global.channels;
}
