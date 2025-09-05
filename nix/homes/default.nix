{
  flake,
  extraArgs,
  nodes,
  extraModules ? []
}:
let
  hmConf =
    {
      modules,
      pkgs,
      extraSpecialArgs ? { },
    }:
    import "${flake.inputs.home-manager}/modules" {
      inherit pkgs;
      extraSpecialArgs = extraArgs // extraSpecialArgs // { inherit pkgs; };
      configuration = {
        imports = modules ++ extraModules;
      };
    };
in
builtins.mapAttrs (k: v: hmConf v) nodes
