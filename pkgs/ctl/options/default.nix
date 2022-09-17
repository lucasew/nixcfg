{ lib, pkgs, ... }:
let
  inherit (builtins) toJSON typeOf tryEval;
  inherit (lib) concatStringsSep attrNames attrValues tail escapeShellArg mapAttrs;
  successOrNull = value: let
    valueDealed = tryEval value;
  in if valueDealed.success
    then valueDealed.value
    else null
  ;

  repr = value: let
    valueDealed = tryEval (value null);
    type = typeOf valueDealed.value;
  in if valueDealed.success
    then {
      function = "<function>";
      string = toString valueDealed.value;
      int = toString valueDealed.value;
    }.${type} or "<unimplemented>"
    else "<failed>";

  recurseDoc = {option, name ? []}: successOrNull (let
    isEndOfLine = (option._type or null) != null;
  in if isEndOfLine
    then {
      action.bash = ''
        echo "${concatStringsSep "." name}"
        echo ${escapeShellArg (option.description or "")}
        echo Default: ${escapeShellArg (repr (v: option.defaultText or option.default or null))}
        echo Example: ${escapeShellArg (repr (v: option.example or null))}
        echo Locations: ${escapeShellArg (concatStringsSep " " option.loc)}
      '';
        # ${concatStringsSep "\n" (map (item: "echo ${escapeShellArg (toString item)}") (attrNames option))}
    } else {
      # action.bash = concatStringsSep "\n" (attrValues (mapAttrs (k: v:
      #   ''
      #     ${concatStringsSep "\n" (map (item: "echo ${escapeShellArg item}") (attrNames option))}
      #   ''
      # ) option));
      subcommands = mapAttrs (k: v: recurseDoc { option = v; name = name ++ [ k ]; }) option;
    });
in {
  subcommands.options = {
    action.bash = ''
      echo "Use --help to see the available options";
    '';
    subcommands = mapAttrs (k: v: recurseDoc { option = v.options; name = [ k ];}) {
      inherit (pkgs.flake.outputs.nixosConfigurations) riverwood whiterun;
    };
  };
}
