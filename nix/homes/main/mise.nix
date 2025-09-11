{
  lib,
  ...
}:

let
  standard-config = ../../../mise.home.toml;
in

{
  config = {
    xdg.configFile."mise/conf.d/00-home-manager.toml" = {
      enable = true;
      source = standard-config;      
    };
  };
}
