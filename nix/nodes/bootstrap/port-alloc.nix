{ config, lib, ... }:

let
  inherit (builtins) removeAttrs;
  inherit (lib) mkOption types submodule literalExpression mdDoc mkDefault attrNames foldl' mapAttrs mkEnableOption attrValues;
in

{
  options.networking.ports = mkOption {
    default = {};

    example = literalExpression ''{
      {
        app.enable = true;
      }
    }'';

    description = "Build time port allocations for services that are only used internally";

    type = types.attrsOf (types.submodule ({ name, config, options, ... }: {
      options = {
        enable = mkEnableOption "Enable automatic port allocation for service ${name}";
        port = mkOption {
          description = "Allocated port for service ${name}";
          type = types.nullOr types.port;
          default = null;
        };
      };
      }));
  };

  imports = [
    ({config, ...}: {
      imports = lib.pipe config.networking.ports [
        (attrNames) # gets only the names of the ports
        (foldl' (x: y: x // { "${y}" = x._port; _port = x._port - 1; })  {_port = 65535; }) # gets the count down of the ports
        (x: removeAttrs x ["_port"]) # removes the utility _port entity
        (mapAttrs (k: v: {config, ...}: { # generates a module for each port but in an attrset
            config.networking.ports.${k}.port = mkDefault v;
        }))
        (attrValues) # convert that module attrset to a list of modules
      ];
    })
  ];
}
