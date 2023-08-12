{ config, lib, ... }:

let
  inherit (builtins) removeAttrs;
  inherit (lib) mkOption types submodule literalExpression mdDoc mkDefault attrNames foldl' mapAttrs mkEnableOption attrValues isString length head warn optional concatStringsSep filter;

  mkAllocModule = {
    description ? null,
    enableDescription ? "Enable allocation of * undefined item *",
    valueKey ? "value",
    valueType ? null,
    valueApply ? null,
    valueLiteral ? value: ''"${keyFunc value}"'',
    valueDescription ? "Allocated value for * undefined item *",
    firstValue ? 0, # first item will get this value allocated
    keyFunc ? toString, # how to generate a conflict checking key from value
    succFunc, # how to get the next item from a item, like next port
    validateFunc ? valueType.check or (value: true), # how to validate one value
    cfg,
    keyPath ? "", # like "networking.ports"

    example ? null,
    internal ? null,
    relatedPackages ? null,
    visible ? null,
    dontThrow ? false,
  }: mkOption {
    inherit
      description
      example
      internal
      relatedPackages
      visible
    ;
    default = {};
    apply = _items: let
      items = removeAttrs _items (filter (item: !_items.${item}.enable) (attrNames _items));
      itemKeys = attrNames items;
      isItemDefined = itemKey: items.${itemKey}.${valueKey} != null;
      definedItems = lib.filter isItemDefined itemKeys;
      undefinedItems = lib.filter (x: !isItemDefined x) itemKeys;

      conflictDict = foldl' (x: y: x // (let
        thisItem = items.${y}.${valueKey};
        thisConflictKey = keyFunc thisItem;
        conflictKey = x.${thisConflictKey} or null;
        isConflict = conflictKey != null;
        isValid = validateFunc thisItem;
        hasProblem = isConflict || !isValid;
      in {
        "${thisConflictKey}" = y; # if validateFunc thisItem then y else warn "Found invalid value at key ${y}. Suggestion: change ${valueKey} to `${suggestedValueLiteral}`" y;
        _conflict = if (x._conflict or null) != null then x._conflict else (if isConflict then {from = conflictKey; to = y; } else null);
        _invalid = (x._invalid or []) ++ (optional (!isValid) y);
      })) {} definedItems;

      getFullKey = key: concatStringsSep "." ((optional (keyPath != "") keyPath) ++ [ key ]);

      isValueConflicts = value: (conflictDict.${keyFunc value} or null) != null;
      suggestValue = prevValue: if (isValueConflicts prevValue) || (!(validateFunc prevValue)) then suggestValue (succFunc prevValue) else prevValue;

      suggestedValue = suggestValue firstValue;
      suggestedValueLiteral = valueLiteral suggestedValue;

      throwPassthru = message: passthru: throw message;

      handleCondition = fallback: condition: message: passthru:
      if condition then passthru
      else if dontThrow
        then
          (passthru // {_message = message; _conflictDict = conflictDict; })
        else
          fallback message passthru;

      handleConditionThrow = condition: message: passthru:
        handleCondition throwPassthru condition message passthru;

      handleConditionWarn = condition: message: passthru:
        handleCondition lib.warn      condition message passthru;

      handleMissingKeyPath = handleConditionWarn (keyPath != "") "mkAllocModule: keyPath missing. Error messages will be less useful";
      handleMissingValues = handleConditionThrow (length undefinedItems == 0) "Key ${getFullKey (head undefinedItems)} is missing a value. Suggestion: set the value to: `${suggestedValueLiteral}`";
      handleConflicts = passthru: handleConditionThrow (conflictDict._conflict == null) "Key ${getFullKey conflictDict._conflict.from} and ${getFullKey conflictDict._conflict.to} have the same values. Suggestion: change the value of one of them to: `${suggestedValueLiteral}`";
      handleInvalidValues = passthru: handleConditionThrow (length conflictDict._invalid == 0) "The following keys have invalid values: ${concatStringsSep ", "(map (getFullKey) conflictDict._invalid)}. Suggestion: change the value of the first key to: `${suggestedValueLiteral}`";

    in lib.pipe items [
      handleMissingKeyPath
      handleMissingValues
      handleConflicts
      handleInvalidValues
    ];

    type = types.attrsOf (types.submodule ({ name, config, options, ... }: {
      options = {
        enable = mkEnableOption (if isString enableDescription then enableDescription else enableDescription name);

        "${valueKey}" = mkOption {
          default = null;

          description = mkEnableOption (if isString valueDescription then valueDescription else valueDescription name);
          type = types.nullOr valueType;
        };
      };
    }));
  };

in {
  options.networking.ports = mkAllocModule {
    dontThrow = true;

    valueKey = "port";
    valueType = types.port;
    cfg = config.networking.ports;
    description = "Build time port allocations for services that are only used internally";
    enableDescription = name: "Enable automatic port allocation for service ${name}";
    valueDescription = name: "Allocated port for service ${name}";

    firstValue = 49151;
    succFunc = x: x - 1;
    valueLiteral = toString;
    validateFunc = x: (types.port.check x) && (x > 1024);
    keyPath = "networking.ports";
    example = literalExpression ''{
      app = {
        enable = true;
        port = 42069; # guided
      };
    }'';
  };

  config.environment.etc = lib.pipe config.networking.ports [
    (attrNames)
    (foldl' (x: y: x // {
      "ports/${y}" = {
        inherit (config.networking.ports.${y}) enable;
        text = toString config.networking.ports.${y}.port;
      };
    }) {})
  ];
  config.networking.ports = {
    eoq = {
      enable = false;
      port = 22;
    };
    trabson = {
      enable = true;
      port = 49139;
    };
  };
}
