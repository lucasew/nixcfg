path:
let
    lsName = import ./lsName.nix;
    importItem = item: import (fn item);
in map importItem (lsName path)