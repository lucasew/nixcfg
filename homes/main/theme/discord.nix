{ pkgs, ... }:
{
  home.file.".config/BetterDiscord/data/stable/custom.css".text = builtins.readFile ./discord.css;
}
