path:
let
    lsName = import ./lsName.nix;
    fn = k: path + "/${k}";
    importItem = item: import (fn item);
in map importItem (lsName path)